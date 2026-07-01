
package repository

import "github.com/kaw3-q/api-blog-go/internal/models"

type PostRepository interface {
	GetAll() []models.Post
	GetByID(id uint) (models.Post, error)
	Create(post models.Post) models.Post
}

type UserRepository interface {
	GetByEmail(email string) (models.User, error)
	GetByID(id uint) (models.User, error)
	Create(user models.User) models.User
	GetAll() []models.User
}
