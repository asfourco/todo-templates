package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/asfourco/todo-templates/backend/db"
	"github.com/asfourco/todo-templates/backend/models"
	"github.com/streamingfast/dhttp"
	"github.com/streamingfast/validator"
)

type TodoController struct {
}

func (tc *TodoController) CreateTodo(dbClient *db.PostgresClient, r *http.Request) (interface{}, error) {
	request := models.CreateTodoRequest{}
	err := dhttp.ExtractJSONRequest(r.Context(), r, &request, dhttp.NewRequestValidator(validator.Rules{
		"Title":       []string{"required"},
		"Description": []string{},
	}))
	if err != nil {
		return nil, err
	}

	var columns []string
	var args []interface{}

	if request.Title != "" {
		columns = append(columns, "Title")
		args = append(args, request.Title)
	}

	zlog.Info("CreateTodo", zap.Any("request", request))
	todo, err := dbClient.Insert("todos", columns, args)
	if err != nil {
		zlog.Error("Failed to insert todo", zap.Error(err))
		return nil, err
	}
	var data models.Todo
	err = todo.Scan(&data.Id, &data.Title, &data.Active, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		zlog.Error("Failed to scan todo", zap.Error(err))
		return nil, err
	}

	return &models.CreateTodoResponse{Todo: data}, nil
}

func (tc *TodoController) GetTodo(dbClient *db.PostgresClient, r *http.Request) (interface{}, error) {
	request := models.GetTodoRequest{}
	err := dhttp.ExtractJSONRequest(r.Context(), r, &request, dhttp.NewRequestValidator(validator.Rules{
		"Id": []string{"required"},
	}))
	if err != nil {
		return nil, err
	}

	row := dbClient.SelectOne("todos", "*", "id = "+strconv.Itoa(request.Id))
	var data models.Todo
	err = row.Scan(&data.Id, &data.Title, &data.Active, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		zlog.Error("Failed to scan todo", zap.Error(err))
	}

	return &models.GetTodoResponse{Todo: data}, nil
}

func (tc *TodoController) GetTodoList(dbClient *db.PostgresClient, r *http.Request) (interface{}, error) {
	request := models.GetTodosRequest{Page: 1, PageSize: db.DEFAULT_PAGE_SIZE}
	err := dhttp.ExtractJSONRequest(r.Context(), r, &request, dhttp.NewRequestValidator(validator.Rules{
		"Page":      []string{"required"},
		"PageSize":  []string{"required"},
		"Completed": []string{},
	}))
	if err != nil {
		return nil, err
	}

	rows, err := dbClient.SelectAll("todos", "id, complete, title, description, created_at", request.Page, request.PageSize)
	if err != nil {
		zlog.Error("Failed to select todos", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var data []models.Todo
	for rows.Next() {
		var todo models.Todo
		err = rows.Scan(&todo.Id, &todo.Active, &todo.Title, &todo.CreatedAt)
		if err != nil {
			zlog.Error("Failed to scan todo", zap.Error(err))
		} else {
			data = append(data, todo)
		}
	}
	return &models.GetTodosResponse{Page: request.Page, ItemsPerPage: request.PageSize, Items: data}, nil
}

func (tc *TodoController) UpdateTodo(dbClient *db.PostgresClient, r *http.Request) (interface{}, error) {
	request := models.UpdateTodoRequest{}
	err := dhttp.ExtractJSONRequest(r.Context(), r, &request, dhttp.NewRequestValidator(validator.Rules{
		"id":          []string{"required"},
		"title":       []string{},
		"description": []string{},
		"active":      []string{},
	}))
	if err != nil {
		return nil, err
	}

	var updates []string
	var args []interface{}

	if request.Title != "" {
		updates = append(updates, "title = $"+strconv.Itoa(len(args)+1))
		args = append(args, request.Title)
	}

	if request.Active != nil {
		updates = append(updates, "complete = $"+strconv.Itoa(len(args)+1))
		args = append(args, request.Active)
	}

	if len(args) == 0 {
		return nil, errors.New("no fields to update")
	}

	updateString := strings.Join(updates, ", ")
	condition := fmt.Sprintf("id = %d", request.Id)

	err = dbClient.Update("todos", updateString, condition, args)
	if err != nil {
		zlog.Error("Failed to update todo", zap.Error(err))
		return nil, err
	}
	return &models.UpdateTodoResponse{}, nil
}

func (tc *TodoController) DeleteTodo(dbClient *db.PostgresClient, r *http.Request) (interface{}, error) {
	request := models.DeleteTodoRequest{}
	err := dhttp.ExtractJSONRequest(r.Context(), r, &request, dhttp.NewRequestValidator(validator.Rules{
		"Id": []string{"required"},
	}))
	if err != nil {
		return nil, err
	}

	condition := fmt.Sprintf("id = %d", request.Id)

	err = dbClient.Delete("todos", condition)
	if err != nil {
		zlog.Error("Failed to delete todo", zap.Error(err))
		return &models.DeleteTodoResponse{Deleted: false}, err
	}
	return &models.DeleteTodoResponse{Deleted: true}, nil
}
