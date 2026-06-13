package repository

import (
	"github.com/kaw3-q/api-blog-go/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	GetAll() []models.Post
	GetByID(id uint) (models.Post, error)
	Create(post models.Post) models.Post
}

type sqlPostRepository struct {
	db *gorm.DB
}

func NewSQLPostRepository(db *gorm.DB) PostRepository {
	return &sqlPostRepository{db: db}
}

func (r *sqlPostRepository) GetAll() []models.Post {
	var posts []models.Post
	r.db.Find(&posts)
	return posts
}

func (r *sqlPostRepository) GetByID(id uint) (models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, id).Error; err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (r *sqlPostRepository) Create(post models.Post) models.Post {
	// GORM usa prepared statements por padrão, protegendo contra SQL Injection
	r.db.Create(&post)
	return post
}
