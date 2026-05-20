//go:build integration

package user_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/user"
	userrepo "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/repository/user"
)

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("rbk_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "postgres start: %v\n", err)
		os.Exit(1)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		fmt.Fprintf(os.Stderr, "dsn: %v\n", err)
		os.Exit(1)
	}

	if err := migrate(dsn); err != nil {
		_ = container.Terminate(ctx)
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}

	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		_ = container.Terminate(ctx)
		fmt.Fprintf(os.Stderr, "pool: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	pool.Close()
	_ = container.Terminate(ctx)
	os.Exit(code)
}

func migrate(dsn string) error {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(sqlDB, migrationsDir())
}

func migrationsDir() string {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "database", "migrations")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("go.mod not found")
		}
		dir = parent
	}
}

func reset(t *testing.T) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		"TRUNCATE users, user_cities, weather_history RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func newRepo(t *testing.T) (context.Context, *userrepo.Repository) {
	t.Helper()
	reset(t)
	return context.Background(), userrepo.NewRepository(pool)
}

func makeUser(t *testing.T, email string) (user.User, user.Password) {
	t.Helper()
	pwd, err := user.NewPassword("secret-1234")
	require.NoError(t, err)
	return user.User{
		ID:        uuid.New(),
		FirstName: "Daniil",
		LastName:  "Kalts",
		Email:     email,
		Role:      user.RoleUser,
	}, pwd
}

func seed(t *testing.T, ctx context.Context, repo *userrepo.Repository, email string) *user.User {
	t.Helper()
	u, pwd := makeUser(t, email)
	created, err := repo.Create(ctx, u, pwd)
	require.NoError(t, err)
	return created
}

func TestUserRepository_Create(t *testing.T) {
	t.Run("success persists user with timestamps", func(t *testing.T) {
		ctx, repo := newRepo(t)
		u, pwd := makeUser(t, "daniil@example.com")

		got, err := repo.Create(ctx, u, pwd)

		require.NoError(t, err)
		assert.Equal(t, u.ID, got.ID)
		assert.Equal(t, "Daniil", got.FirstName)
		assert.Equal(t, "daniil@example.com", got.Email)
		assert.Equal(t, user.RoleUser, got.Role)
		assert.False(t, got.CreatedAt.IsZero())
		assert.False(t, got.UpdatedAt.IsZero())
	})

	t.Run("duplicate email returns ErrEmailAlreadyExists", func(t *testing.T) {
		ctx, repo := newRepo(t)
		seed(t, ctx, repo, "dup@example.com")

		u, pwd := makeUser(t, "dup@example.com")
		_, err := repo.Create(ctx, u, pwd)

		assert.ErrorIs(t, err, user.ErrEmailAlreadyExists)
	})

	t.Run("duplicate email after soft-delete is allowed", func(t *testing.T) {
		ctx, repo := newRepo(t)
		first := seed(t, ctx, repo, "reuse@example.com")
		require.NoError(t, repo.SoftDelete(ctx, first.ID))

		u, pwd := makeUser(t, "reuse@example.com")
		got, err := repo.Create(ctx, u, pwd)

		require.NoError(t, err)
		assert.NotEqual(t, first.ID, got.ID)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, repo := newRepo(t)
		created := seed(t, ctx, repo, "g1@example.com")

		got, err := repo.GetByID(ctx, created.ID)

		require.NoError(t, err)
		assert.Equal(t, created.ID, got.ID)
		assert.Equal(t, created.Email, got.Email)
	})

	t.Run("not found", func(t *testing.T) {
		ctx, repo := newRepo(t)
		_, err := repo.GetByID(ctx, uuid.New())
		assert.ErrorIs(t, err, user.ErrNotFound)
	})

	t.Run("soft-deleted user is hidden", func(t *testing.T) {
		ctx, repo := newRepo(t)
		created := seed(t, ctx, repo, "g2@example.com")
		require.NoError(t, repo.SoftDelete(ctx, created.ID))

		_, err := repo.GetByID(ctx, created.ID)
		assert.ErrorIs(t, err, user.ErrNotFound)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, repo := newRepo(t)
		created := seed(t, ctx, repo, "ge@example.com")

		got, err := repo.GetByEmail(ctx, "ge@example.com")

		require.NoError(t, err)
		assert.Equal(t, created.ID, got.ID)
	})

	t.Run("not found", func(t *testing.T) {
		ctx, repo := newRepo(t)
		_, err := repo.GetByEmail(ctx, "missing@example.com")
		assert.ErrorIs(t, err, user.ErrNotFound)
	})

	t.Run("soft-deleted user is hidden", func(t *testing.T) {
		ctx, repo := newRepo(t)
		created := seed(t, ctx, repo, "hide@example.com")
		require.NoError(t, repo.SoftDelete(ctx, created.ID))

		_, err := repo.GetByEmail(ctx, "hide@example.com")
		assert.ErrorIs(t, err, user.ErrNotFound)
	})
}

func TestUserRepository_GetCredentialsByEmail(t *testing.T) {
	t.Run("success returns password hash and salt", func(t *testing.T) {
		ctx, repo := newRepo(t)
		seed(t, ctx, repo, "cred@example.com")

		u, pwd, err := repo.GetCredentialsByEmail(ctx, "cred@example.com")

		require.NoError(t, err)
		assert.Equal(t, "cred@example.com", u.Email)
		assert.NotEmpty(t, pwd.Hash)
		assert.NotEmpty(t, pwd.Salt)
		assert.True(t, pwd.Matches("secret-1234"))
	})

	t.Run("not found", func(t *testing.T) {
		ctx, repo := newRepo(t)
		_, _, err := repo.GetCredentialsByEmail(ctx, "missing@example.com")
		assert.ErrorIs(t, err, user.ErrNotFound)
	})

	t.Run("soft-deleted user cannot authenticate", func(t *testing.T) {
		ctx, repo := newRepo(t)
		created := seed(t, ctx, repo, "deleted@example.com")
		require.NoError(t, repo.SoftDelete(ctx, created.ID))

		_, _, err := repo.GetCredentialsByEmail(ctx, "deleted@example.com")
		assert.ErrorIs(t, err, user.ErrNotFound)
	})
}

func TestUserRepository_List(t *testing.T) {
	t.Run("empty returns empty slice", func(t *testing.T) {
		ctx, repo := newRepo(t)
		got, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("returns only non-deleted users", func(t *testing.T) {
		ctx, repo := newRepo(t)
		keep := seed(t, ctx, repo, "keep@example.com")
		gone := seed(t, ctx, repo, "gone@example.com")
		require.NoError(t, repo.SoftDelete(ctx, gone.ID))

		got, err := repo.List(ctx)

		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, keep.ID, got[0].ID)
	})
}

func TestUserRepository_Update(t *testing.T) {
	t.Run("success bumps updated_at and persists fields", func(t *testing.T) {
		ctx, repo := newRepo(t)
		created := seed(t, ctx, repo, "upd@example.com")

		got, err := repo.Update(ctx, user.User{
			ID:        created.ID,
			FirstName: "New",
			LastName:  "Name",
			Email:     "new@example.com",
		})

		require.NoError(t, err)
		assert.Equal(t, "New", got.FirstName)
		assert.Equal(t, "Name", got.LastName)
		assert.Equal(t, "new@example.com", got.Email)
		assert.True(t, got.UpdatedAt.After(created.UpdatedAt) || got.UpdatedAt.Equal(created.UpdatedAt))
	})

	t.Run("not found", func(t *testing.T) {
		ctx, repo := newRepo(t)
		_, err := repo.Update(ctx, user.User{
			ID:        uuid.New(),
			FirstName: "X",
			LastName:  "Y",
			Email:     "x@example.com",
		})
		assert.ErrorIs(t, err, user.ErrNotFound)
	})

	t.Run("soft-deleted user is not updatable", func(t *testing.T) {
		ctx, repo := newRepo(t)
		created := seed(t, ctx, repo, "sd@example.com")
		require.NoError(t, repo.SoftDelete(ctx, created.ID))

		_, err := repo.Update(ctx, user.User{
			ID:        created.ID,
			FirstName: "X",
			LastName:  "Y",
			Email:     "z@example.com",
		})
		assert.ErrorIs(t, err, user.ErrNotFound)
	})

	t.Run("duplicate email returns ErrEmailAlreadyExists", func(t *testing.T) {
		ctx, repo := newRepo(t)
		taken := seed(t, ctx, repo, "taken@example.com")
		mover := seed(t, ctx, repo, "mover@example.com")

		_, err := repo.Update(ctx, user.User{
			ID:        mover.ID,
			FirstName: mover.FirstName,
			LastName:  mover.LastName,
			Email:     taken.Email,
		})
		assert.ErrorIs(t, err, user.ErrEmailAlreadyExists)
	})
}

func TestUserRepository_SoftDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, repo := newRepo(t)
		created := seed(t, ctx, repo, "sd1@example.com")

		err := repo.SoftDelete(ctx, created.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, created.ID)
		assert.ErrorIs(t, err, user.ErrNotFound)
	})

	t.Run("idempotent second call returns ErrNotFound", func(t *testing.T) {
		ctx, repo := newRepo(t)
		created := seed(t, ctx, repo, "sd2@example.com")
		require.NoError(t, repo.SoftDelete(ctx, created.ID))

		err := repo.SoftDelete(ctx, created.ID)
		assert.ErrorIs(t, err, user.ErrNotFound)
	})

	t.Run("unknown id returns ErrNotFound", func(t *testing.T) {
		ctx, repo := newRepo(t)
		err := repo.SoftDelete(ctx, uuid.New())
		assert.ErrorIs(t, err, user.ErrNotFound)
	})
}

func TestUserRepository_GetByID_WrapsUnexpectedErrors(t *testing.T) {
	ctx, repo := newRepo(t)
	cancelled, cancel := context.WithCancel(ctx)
	cancel()

	_, err := repo.GetByID(cancelled, uuid.New())

	require.Error(t, err)
	assert.False(t, errors.Is(err, user.ErrNotFound))
	assert.Contains(t, fmt.Sprint(err), "получение пользователя по id")
}
