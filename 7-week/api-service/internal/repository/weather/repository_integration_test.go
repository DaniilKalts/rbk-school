//go:build integration

package weather_test

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

	"github.com/DaniilKalts/rbk-school/7-week/internal/domain/history"
	"github.com/DaniilKalts/rbk-school/7-week/internal/domain/user"

	userrepo "github.com/DaniilKalts/rbk-school/7-week/internal/repository/user"
	weatherrepo "github.com/DaniilKalts/rbk-school/7-week/internal/repository/weather"
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

func setup(t *testing.T) (context.Context, *weatherrepo.Repository, *user.User) {
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

	return ctx, weatherrepo.NewRepository(pool), u
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

func newHistory(userID uuid.UUID, city string, temp float64) history.History {
	return history.History{
		ID:          uuid.New(),
		UserID:      userID,
		City:        city,
		Temperature: temp,
		Description: "clear sky",
	}
}

func insert(t *testing.T, ctx context.Context, repo *weatherrepo.Repository, h history.History) *history.History {
	t.Helper()
	got, err := repo.CreateHistory(ctx, h)
	require.NoError(t, err)
	time.Sleep(2 * time.Millisecond)
	return got
}

func TestWeatherRepository_CreateHistory(t *testing.T) {
	t.Run("success persists row", func(t *testing.T) {
		ctx, repo, u := setup(t)

		got, err := repo.CreateHistory(ctx, newHistory(u.ID, "Almaty", 12.34))

		require.NoError(t, err)
		assert.Equal(t, u.ID, got.UserID)
		assert.Equal(t, "Almaty", got.City)
		assert.InDelta(t, 12.34, got.Temperature, 0.001)
		assert.Equal(t, "clear sky", got.Description)
		assert.False(t, got.RequestedAt.IsZero())
	})

	t.Run("unknown user id violates FK", func(t *testing.T) {
		ctx, repo, _ := setup(t)
		_, err := repo.CreateHistory(ctx, newHistory(uuid.New(), "Almaty", 1.0))
		require.Error(t, err)
	})
}

func TestWeatherRepository_ListHistory(t *testing.T) {
	t.Run("empty returns empty slice", func(t *testing.T) {
		ctx, repo, u := setup(t)
		got, err := repo.ListHistory(ctx, u.ID, "", 0, 0)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("isolated per user", func(t *testing.T) {
		ctx, repo, u := setup(t)
		other := seedUser(t, ctx)

		insert(t, ctx, repo, newHistory(u.ID, "Almaty", 1))
		insert(t, ctx, repo, newHistory(other.ID, "Almaty", 2))

		got, err := repo.ListHistory(ctx, u.ID, "", 0, 0)

		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, u.ID, got[0].UserID)
	})

	t.Run("city filter matches exact string", func(t *testing.T) {
		ctx, repo, u := setup(t)
		insert(t, ctx, repo, newHistory(u.ID, "Almaty", 1))
		insert(t, ctx, repo, newHistory(u.ID, "Astana", 2))

		got, err := repo.ListHistory(ctx, u.ID, "Almaty", 0, 0)

		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, "Almaty", got[0].City)
	})

	t.Run("empty city filter returns all", func(t *testing.T) {
		ctx, repo, u := setup(t)
		insert(t, ctx, repo, newHistory(u.ID, "Almaty", 1))
		insert(t, ctx, repo, newHistory(u.ID, "Astana", 2))

		got, err := repo.ListHistory(ctx, u.ID, "", 0, 0)

		require.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("ordered by requested_at DESC", func(t *testing.T) {
		ctx, repo, u := setup(t)
		first := insert(t, ctx, repo, newHistory(u.ID, "Almaty", 1))
		second := insert(t, ctx, repo, newHistory(u.ID, "Almaty", 2))
		third := insert(t, ctx, repo, newHistory(u.ID, "Almaty", 3))

		got, err := repo.ListHistory(ctx, u.ID, "", 0, 0)

		require.NoError(t, err)
		require.Len(t, got, 3)
		assert.Equal(t, third.ID, got[0].ID)
		assert.Equal(t, second.ID, got[1].ID)
		assert.Equal(t, first.ID, got[2].ID)
	})

	t.Run("limit caps result count", func(t *testing.T) {
		ctx, repo, u := setup(t)
		insert(t, ctx, repo, newHistory(u.ID, "Almaty", 1))
		insert(t, ctx, repo, newHistory(u.ID, "Almaty", 2))
		insert(t, ctx, repo, newHistory(u.ID, "Almaty", 3))

		got, err := repo.ListHistory(ctx, u.ID, "", 2, 0)

		require.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("limit 0 means unlimited", func(t *testing.T) {
		ctx, repo, u := setup(t)
		for i := 0; i < 3; i++ {
			insert(t, ctx, repo, newHistory(u.ID, "Almaty", float64(i)))
		}

		got, err := repo.ListHistory(ctx, u.ID, "", 0, 0)

		require.NoError(t, err)
		assert.Len(t, got, 3)
	})

	t.Run("offset skips leading rows", func(t *testing.T) {
		ctx, repo, u := setup(t)
		first := insert(t, ctx, repo, newHistory(u.ID, "Almaty", 1))
		_ = insert(t, ctx, repo, newHistory(u.ID, "Almaty", 2))
		third := insert(t, ctx, repo, newHistory(u.ID, "Almaty", 3))

		got, err := repo.ListHistory(ctx, u.ID, "", 0, 1)

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.NotEqual(t, third.ID, got[0].ID)
		assert.Equal(t, first.ID, got[1].ID)
	})

	t.Run("offset past end returns empty", func(t *testing.T) {
		ctx, repo, u := setup(t)
		insert(t, ctx, repo, newHistory(u.ID, "Almaty", 1))

		got, err := repo.ListHistory(ctx, u.ID, "", 0, 10)

		require.NoError(t, err)
		assert.Empty(t, got)
	})
}

func TestWeatherRepository_Cascade(t *testing.T) {
	t.Run("deleting user wipes their history", func(t *testing.T) {
		ctx, repo, u := setup(t)
		insert(t, ctx, repo, newHistory(u.ID, "Almaty", 1))

		_, err := pool.Exec(ctx, "DELETE FROM users WHERE id = $1", u.ID)
		require.NoError(t, err)

		got, err := repo.ListHistory(ctx, u.ID, "", 0, 0)
		require.NoError(t, err)
		assert.Empty(t, got)
	})
}
