package api

import (
	"net/http"

	"github.com/asfourco/templates/backend/controllers"
)

func (s *Server) GetTodoList(r *http.Request) (interface{}, error) {
	tc := controllers.TodoController{}
	result, err := tc.GetTodoList(s.postgresClient, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) GetTodo(r *http.Request) (interface{}, error) {
	tc := controllers.TodoController{}
	result, err := tc.GetTodo(s.postgresClient, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) CreateTodo(r *http.Request) (interface{}, error) {
	tc := controllers.TodoController{}
	result, err := tc.CreateTodo(s.postgresClient, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) UpdateTodo(r *http.Request) (interface{}, error) {
	tc := controllers.TodoController{}
	result, err := tc.UpdateTodo(s.postgresClient, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) DeleteTodo(r *http.Request) (interface{}, error) {
	tc := controllers.TodoController{}
	result, err := tc.DeleteTodo(s.postgresClient, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}
