package models

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Status bool   `json:"status"`
	UserID string `json:"user_id"`
}
