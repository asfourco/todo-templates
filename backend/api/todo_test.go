package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/asfourco/todo-templates/backend/models"
	"strings"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/asfourco/todo-templates/backend/api"
	"github.com/asfourco/todo-templates/backend/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodo(t *testing.T) {
	dbUrl := setupTestDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	postgresClient, err := db.NewPostgresClient(ctx, dbUrl)
	if err != nil {
		t.Fatal(err)
	}
	s, err := api.NewServer(ctx, "8080", postgresClient)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name               string
		body               string
		wantErr            bool
		expectedStatusCode int
	}{
		{
			name:               "valid request",
			body:               `{"title": "test todo"}`,
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "invalid request",
			body:               `{}`,
			wantErr:            true,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/todos", strings.NewReader(tt.body))
			assert.NoError(t, err, "error creating test request")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp, err := s.CreateTodo(r)
				if (err != nil) != tt.wantErr {
					t.Errorf("CreateTodoRequest() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err == nil {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(resp)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			})

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestGetTodoById(t *testing.T) {
	dbUrl := setupTestDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	postgresClient, err := db.NewPostgresClient(ctx, dbUrl)
	if err != nil {
		t.Fatal(err)
	}
	setupTodos(t, postgresClient)
	s, err := api.NewServer(ctx, "8080", postgresClient)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "valid request",
			id:      1,
			wantErr: false,
		},
		{
			name:    "invalid request",
			id:      -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/todos/?id=%d", tt.id), nil)
			assert.NoError(t, err, "error creating test request")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp, err := s.GetTodo(r)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetTodo() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err == nil {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(resp)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			})

			handler.ServeHTTP(rr, req)

			if tt.wantErr {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			} else {
				assert.Equal(t, http.StatusOK, rr.Code)
			}

		})
	}

}

func TestGetTodoList(t *testing.T) {
	dbUrl := setupTestDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	postgresClient, err := db.NewPostgresClient(ctx, dbUrl)
	if err != nil {
		t.Fatal(err)
	}
	setupTodos(t, postgresClient)
	s, err := api.NewServer(ctx, "8080", postgresClient)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		page           uint32
		pageSize       uint32
		expectedLength int
	}{
		{
			name:           "valid request",
			page:           0,
			pageSize:       10,
			expectedLength: 5,
		},
		{
			name:           "reduced page size",
			page:           0,
			pageSize:       2,
			expectedLength: 2,
		},
		{
			name:           "invalid page",
			page:           1,
			pageSize:       10,
			expectedLength: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/todos/?page=%d&page_size=%d", test.page, test.pageSize), nil)
			assert.NoError(t, err, "error creating test request")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp, err := s.GetTodoList(r)
				if err != nil {
					t.Errorf("GetTodoList() error = %v", err)
					return
				}

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			})

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			var resp models.GetTodosResponse
			err = json.NewDecoder(rr.Body).Decode(&resp)
			assert.NoError(t, err, "error decoding response")

			assert.Equal(t, test.expectedLength, len(resp.Items))
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	dbUrl := setupTestDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	postgresClient, err := db.NewPostgresClient(ctx, dbUrl)
	if err != nil {
		t.Fatal(err)
	}
	setupTodos(t, postgresClient)
	s, err := api.NewServer(ctx, "8080", postgresClient)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name               string
		body               string
		wantErr            bool
		expectedStatusCode int
	}{
		{
			name:               "valid request",
			body:               `{"id": 1, "title": "test todo", "active": false}`,
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "invalid request",
			body:               `{}`,
			wantErr:            true,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/todos", strings.NewReader(tt.body))
			assert.NoError(t, err, "error creating test request")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp, err := s.UpdateTodo(r)
				if (err != nil) != tt.wantErr {
					t.Errorf("UpdateTodoRequest() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err == nil {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(resp)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			})

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	dbUrl := setupTestDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	postgresClient, err := db.NewPostgresClient(ctx, dbUrl)
	if err != nil {
		t.Fatal(err)
	}
	setupTodos(t, postgresClient)
	s, err := api.NewServer(ctx, "8080", postgresClient)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name               string
		id                 int
		wantErr            bool
		expectedStatusCode int
	}{
		{
			name:               "valid request",
			id:                 1,
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "invalid request",
			id:                 -1,
			wantErr:            true,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("/todos?id=%d", tt.id), nil)
			assert.NoError(t, err, "error creating test request")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp, err := s.DeleteTodo(r)
				if (err != nil) != tt.wantErr {
					t.Errorf("DeleteTodoRequest() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err == nil {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(resp)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			})

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
