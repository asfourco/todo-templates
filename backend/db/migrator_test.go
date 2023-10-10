package db_test

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/lestrrat-go/backoff"

	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/stretchr/testify/assert"

	"github.com/asfourco/todo-templates/backend/db"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbUser     = "testuser"
	dbPassword = "testpass"
	dbName     = "testdb"
	dbImage    = "docker.io/postgres:alpine"
)

//go:embed migrations/*
var migrationFiles embed.FS

// setupTestDB sets up a test PostgreSQL container.
func setupTestDB(t *testing.T) (string, func()) {
	err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	if err != nil {
		t.Fatalf("couldn't set env: %v", err)
	}

	ctx := context.Background()

	pgC, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage(dbImage),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		if err := pgC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	dbUrl, err := pgC.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	if err != nil {
		t.Fatal(err)
	}
	return dbUrl, cleanup
}

func TestMigration_Migrate(t *testing.T) {
	dbUrl, cleanup := setupTestDB(t)
	defer cleanup()

	sourceDriver, err := httpfs.New(http.FS(migrationFiles), "migrations")
	if err != nil {
		t.Fatal(err)
	}

	m := db.NewMigration(
		dbUrl,
		"",
		sourceDriver,
	)

	backoffPolicy := backoff.NewExponential(
		backoff.WithInterval(time.Second),
		backoff.WithMaxInterval(2*time.Second),
		backoff.WithJitterFactor(0.05),
	)
	b, cancel := backoffPolicy.Start(context.Background())
	defer cancel()

	for backoff.Continue(b) {
		err = m.Migrate()
		if err != nil && err != migrate.ErrNoChange {
			// If error persists after retries, fail
			if !backoff.Continue(b) {
				t.Fatalf("an error occurred while syncing the database schema: %v", err)
			}
			continue
		} else {
			break
		}
	}

	assert.NoError(t, err, fmt.Sprintf("couldn't migrate: %v", err))
}

func TestMigration_Rollback(t *testing.T) {
	dbUrl, cleanup := setupTestDB(t)
	defer cleanup()

	sourceDriver, err := httpfs.New(http.FS(migrationFiles), "migrations")
	if err != nil {
		t.Fatal(err)
	}

	m := db.NewMigration(
		dbUrl,
		"",
		sourceDriver,
	)

	backoffPolicy := backoff.NewExponential(
		backoff.WithInterval(time.Second),
		backoff.WithMaxInterval(2*time.Second),
		backoff.WithJitterFactor(0.05),
	)
	b, cancel := backoffPolicy.Start(context.Background())
	defer cancel()

	for backoff.Continue(b) {
		err = m.Migrate()
		if err != nil && err != migrate.ErrNoChange {
			// If error persists after retries, fail
			if !backoff.Continue(b) {
				t.Fatalf("an error occurred while syncing the database schema: %v", err)
			}
			continue
		} else {
			break
		}
	}

	for backoff.Continue(b) {
		err = m.Rollback()
		if err != nil && err != migrate.ErrNoChange {
			// If error persists after retries, fail
			if !backoff.Continue(b) {
				t.Fatalf("an error occurred while reverting the database schema: %v", err)
			}
			continue
		} else {
			break
		}
	}

	assert.NoError(t, err, fmt.Sprintf("couldn't rollback: %v", err))
}
