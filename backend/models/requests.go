package models

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type GetTodoRequest struct {
	Id int `json:"id"`
}

type GetTodosRequest struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Active   *bool `json:"active"`
}

type UpdateTodoRequest struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Active *bool  `json:"active"`
}

type DeleteTodoRequest struct {
	Id int `json:"id"`
}
