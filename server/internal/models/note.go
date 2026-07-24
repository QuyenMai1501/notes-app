package models

import "time"

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateNoteInput struct {
	Title   string   `json:"title" binding:"required,min=1"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type UpdateNoteInput struct {
	Title   string   `json:"title" binding:"required,min=1"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}
