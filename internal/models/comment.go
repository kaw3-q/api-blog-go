package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	AuthorID  int       `json:"author_id"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
