package models

type Todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Active    *bool  `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
