package api_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/asfourco/todo-templates/backend/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/lestrrat-go/backoff"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbUser  = "testuser"
	dbPass  = "testpass"
	dbName  = "testdb"
	dbImage = "docker.io/postgres:alpine"
)

func setupTestDB(t *testing.T) string {
	// Disable RYUK (aka reaper)
	err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	if err != nil {
		t.Fatalf("unable to set env variable: %v", err)
	}
	ctx := context.Background()

	pgC, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage(dbImage),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgC.Terminate(ctx); err != nil {
			t.Fatalf("unable to terminate container: %v", err)
		}
	})

	dbUrl, err := pgC.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	if err != nil {
		t.Fatal(err)
	}

	backoffPolicy := backoff.NewExponential(
		backoff.WithInterval(time.Second),
		backoff.WithMaxInterval(2*time.Second),
		backoff.WithJitterFactor(0.05),
	)
	b, cancel := backoffPolicy.Start(context.Background())
	defer cancel()

	migrationPath := "file://../db/migrations"
	m, err := migrate.New(migrationPath, dbUrl)
	if err != nil {
		t.Fatalf("unable to initialize migrations: %v", err)
	}

	for backoff.Continue(b) {
		err := m.Up()
		if err != nil && err != migrate.ErrNoChange {
			zlog.Debug("could not apply migrations, retrying")
			if !backoff.Continue(b) {
				t.Fatalf("Error while syncing the database schema: %v", err)
			}
			continue
		} else {
			break
		}
	}

	return dbUrl
}

func setupTodos(t *testing.T, client *db.PostgresClient) {
	var todoItems = []string{"Buy milk", "Buy eggs", "Buy bread", "Buy butter", "Buy cheese"}
	for i := 1; i <= len(todoItems); i++ {
		_, err := client.Insert("todos", []string{"title", "active"}, []interface{}{todoItems[i-1], true})
		if err != nil {
			t.Fatalf("unable to insert todo: %v", err)
		}
	}
}
