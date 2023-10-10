package api

import (
	"net/http"

	"github.com/asfourco/todo-templates/backend/controllers"
)

var tc = controllers.TodoController{}

func (s *Server) GetTodoList(r *http.Request) (resp interface{}, err error) {
	result, err := tc.GetTodoList(s.postgresClient, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) GetTodo(r *http.Request) (resp interface{}, err error) {
	result, err := tc.GetTodo(s.postgresClient, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) CreateTodo(r *http.Request) (resp interface{}, err error) {
	result, err := tc.CreateTodo(s.postgresClient, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) UpdateTodo(r *http.Request) (resp interface{}, err error) {
	result, err := tc.UpdateTodo(s.postgresClient, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) DeleteTodo(r *http.Request) (resp interface{}, err error) {
	tc := controllers.TodoController{}
	result, err := tc.DeleteTodo(s.postgresClient, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}
