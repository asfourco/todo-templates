package models

type CreateTodoResponse struct {
	Todo Todo `json:"todo"`
}

type GetTodoResponse struct {
	Todo Todo `json:"todo"`
}

type GetTodosResponse struct {
	Page         int    `json:"page"`
	ItemsPerPage int    `json:"page_size"`
	Items        []Todo `json:"items"`
}

type UpdateTodoResponse struct {
	Todo Todo `json:"todo"`
}

type HealthResponse struct {
	Healthy bool `json:"healthy"`
}

type DeleteTodoResponse struct {
	Deleted bool `json:"deleted"`
}
