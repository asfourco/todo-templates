package models

import "time"

type Todo struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Active    *bool     `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
