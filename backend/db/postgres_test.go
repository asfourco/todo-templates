package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/asfourco/todo-templates/backend/db"
	"github.com/stretchr/testify/assert"
)

var postgresClient *db.PostgresClient

func TestPostgresClient(t *testing.T) {
	dbUrl, cleanup := setupTestDB(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var err error
	postgresClient, err = db.NewPostgresClient(ctx, dbUrl)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestCreateTable", testCreate)
	t.Run("TestInsert", testInsert)
	t.Run("TestSelect", testSelect)
	t.Run("TestUpdate", testUpdate)
	t.Run("TestDelete", testDelete)
}

func testCreate(t *testing.T) {
	err := postgresClient.CreateTable(
		dbName, "id SERIAL PRIMARY KEY, name TEXT NOT NULL, active BOOLEAN NOT NULL",
	)
	assert.NoError(t, err)
}

func testInsert(t *testing.T) {
	columns := []string{"name", "active"}
	args := []interface{}{"John Doe", true}
	_, err := postgresClient.Insert(dbName, columns, args)
	assert.NoError(t, err)
}

func testSelect(t *testing.T) {
	rows, err := postgresClient.Select(dbName, "name", "name = 'John Doe'", 0, 100)
	assert.NoError(t, err)

	found := false
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		assert.NoError(t, err)

		if name == "John Doe" {
			found = true
			break
		}
	}

	assert.True(t, found)
}

func testUpdate(t *testing.T) {
	err := postgresClient.Update(dbName, "name = $1", "id = 1", []interface{}{"Jane Doe"})
	assert.NoError(t, err)
}

func testDelete(t *testing.T) {
	err := postgresClient.Delete(dbName, "name = 'Jane Doe'")
	assert.NoError(t, err)

	// Check if the row was deleted
	rows, err := postgresClient.Select(dbName, "name", "name = 'Jane Doe'", 0, 100)
	assert.NoError(t, err)

	found := false
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		assert.NoError(t, err)

		if name == "Jane Doe" {
			found = true
			break
		}
	}

	assert.False(t, found)
}
