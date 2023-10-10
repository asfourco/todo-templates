package models

type CreateTodoRequest struct {
	Title *string `json:"title"`
}

type GetTodosRequest struct {
	Page     int    `schema:"page"`
	PageSize int    `schema:"page_size"`
	Active   *bool  `schema:"active"`
	OrderBy  string `schema:"order_by"`
}

type UpdateTodoRequest struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Active *bool  `json:"active"`
}

type GetTodoRequest struct {
	Id int `schema:"id"`
}

type DeleteTodoRequest struct {
	Id int `schema:"id"`
}
