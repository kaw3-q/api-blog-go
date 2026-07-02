package models

import "time"

type PostStatus string

const (
	PostStatusDraft     PostStatus = "DRAFT"
	PostStatusPublished PostStatus = "PUBLISHED"
	PostStatusArchived  PostStatus = "ARCHIVED"
)

type Post struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Excerpt         *string    `json:"excerpt,omitempty"`
	Content         string     `json:"content"`
	CoverImage      *string    `json:"cover_image,omitempty"`
	Status          PostStatus `json:"status"`
	Featured        bool       `json:"featured"`
	Views           int        `json:"views"`
	ReadingTime     *int       `json:"reading_time,omitempty"`
	MetaTitle       *string    `json:"meta_title,omitempty"`
	MetaDescription *string    `json:"meta_description,omitempty"`
	AuthorID        int        `json:"author_id"`
	CategoryID      *int       `json:"category_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`
}
