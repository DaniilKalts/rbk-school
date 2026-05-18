//go:build integration

package city_test

import (
	"context"
	"database/sql"
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

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"

	cityrepo "github.com/DaniilKalts/rbk-school/6-week/internal/repository/city"
	userrepo "github.com/DaniilKalts/rbk-school/6-week/internal/repository/user"
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

func setup(t *testing.T) (context.Context, *cityrepo.Repository, *user.User) {
	t.Helper()
	reset(t)
	ctx := context.Background()

	pwd, err := user.NewPassword("secret-1234")
	require.NoError(t, err)
	u, err := userrepo.NewRepository(pool).Create(ctx, user.User{
		ID:        uuid.New(),
		FirstName: "Daniil",
		LastName:  "Kalts",
		Email:     uuid.NewString() + "@example.com",
		Role:      user.RoleUser,
	}, pwd)
	require.NoError(t, err)

	return ctx, cityrepo.NewRepository(pool), u
}

func seedUser(t *testing.T, ctx context.Context) *user.User {
	t.Helper()
	pwd, err := user.NewPassword("secret-1234")
	require.NoError(t, err)
	u, err := userrepo.NewRepository(pool).Create(ctx, user.User{
		ID:        uuid.New(),
		FirstName: "Daniil",
		LastName:  "Kalts",
		Email:     "x-" + uuid.NewString() + "@example.com",
		Role:      user.RoleUser,
	}, pwd)
	require.NoError(t, err)
	return u
}

func makeCity(userID uuid.UUID, name string) city.City {
	return city.City{ID: uuid.New(), UserID: userID, Name: name}
}

func TestCityRepository_Create(t *testing.T) {
	t.Run("success persists city", func(t *testing.T) {
		ctx, repo, u := setup(t)

		got, err := repo.Create(ctx, makeCity(u.ID, "Almaty"))

		require.NoError(t, err)
		assert.Equal(t, u.ID, got.UserID)
		assert.Equal(t, "Almaty", got.Name)
		assert.False(t, got.CreatedAt.IsZero())
	})

	t.Run("duplicate (user_id, city) returns ErrAlreadyExists", func(t *testing.T) {
		ctx, repo, u := setup(t)
		_, err := repo.Create(ctx, makeCity(u.ID, "Almaty"))
		require.NoError(t, err)

		_, err = repo.Create(ctx, makeCity(u.ID, "Almaty"))
		assert.ErrorIs(t, err, city.ErrAlreadyExists)
	})

	t.Run("same city for different users is allowed", func(t *testing.T) {
		ctx, repo, u := setup(t)
		other := seedUser(t, ctx)

		_, err := repo.Create(ctx, makeCity(u.ID, "Almaty"))
		require.NoError(t, err)
		_, err = repo.Create(ctx, makeCity(other.ID, "Almaty"))
		require.NoError(t, err)
	})

	t.Run("unknown user id violates FK", func(t *testing.T) {
		ctx, repo, _ := setup(t)
		_, err := repo.Create(ctx, makeCity(uuid.New(), "Almaty"))
		require.Error(t, err)
		assert.NotErrorIs(t, err, city.ErrAlreadyExists)
	})
}

func TestCityRepository_ListByUserID(t *testing.T) {
	t.Run("empty returns empty slice", func(t *testing.T) {
		ctx, repo, u := setup(t)
		got, err := repo.ListByUserID(ctx, u.ID)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("returns only cities for the given user", func(t *testing.T) {
		ctx, repo, u := setup(t)
		other := seedUser(t, ctx)

		_, _ = repo.Create(ctx, makeCity(u.ID, "Almaty"))
		_, _ = repo.Create(ctx, makeCity(u.ID, "Astana"))
		_, _ = repo.Create(ctx, makeCity(other.ID, "Shymkent"))

		got, err := repo.ListByUserID(ctx, u.ID)

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.ElementsMatch(t, []string{"Almaty", "Astana"}, []string{got[0].Name, got[1].Name})
	})
}

func TestCityRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, repo, u := setup(t)
		created, err := repo.Create(ctx, makeCity(u.ID, "Almaty"))
		require.NoError(t, err)

		err = repo.Delete(ctx, u.ID, created.ID)
		require.NoError(t, err)

		got, err := repo.ListByUserID(ctx, u.ID)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("unknown city id returns ErrNotFound", func(t *testing.T) {
		ctx, repo, u := setup(t)
		err := repo.Delete(ctx, u.ID, uuid.New())
		assert.ErrorIs(t, err, city.ErrNotFound)
	})

	t.Run("wrong user cannot delete another user's city", func(t *testing.T) {
		ctx, repo, u := setup(t)
		created, err := repo.Create(ctx, makeCity(u.ID, "Almaty"))
		require.NoError(t, err)

		err = repo.Delete(ctx, uuid.New(), created.ID)
		assert.ErrorIs(t, err, city.ErrNotFound)

		got, err := repo.ListByUserID(ctx, u.ID)
		require.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("cascade: deleting user removes cities", func(t *testing.T) {
		ctx, repo, u := setup(t)
		_, err := repo.Create(ctx, makeCity(u.ID, "Almaty"))
		require.NoError(t, err)

		_, err = pool.Exec(ctx, "DELETE FROM users WHERE id = $1", u.ID)
		require.NoError(t, err)

		got, err := repo.ListByUserID(ctx, u.ID)
		require.NoError(t, err)
		assert.Empty(t, got)
	})
}
