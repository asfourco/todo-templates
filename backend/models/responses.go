package models

type GetTodosResponse struct {
	Page         int    `json:"page"`
	ItemsPerPage int    `json:"page_size"`
	Items        []Todo `json:"items"`
}

type HealthResponse struct {
	Healthy bool `json:"healthy"`
}

type DeleteTodoResponse struct {
	Deleted bool `json:"deleted"`
}
