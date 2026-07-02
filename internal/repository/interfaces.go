package repository

import "github.com/kaw3-q/api-blog-go/internal/models"

type PostRepository interface {
	GetAll() ([]models.Post, error)
	GetByID(id int) (models.Post, error)
	Create(post models.Post) (models.Post, error)
}

type UserRepository interface {
	GetByEmail(email string) (models.User, error)
	GetByID(id int) (models.User, error)
	Create(user models.User) (models.User, error)
	GetAll() ([]models.User, error)
}

type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	GetByID(id int) (models.Category, error)
	Create(category models.Category) (models.Category, error)
}

type TagRepository interface {
	GetAll() ([]models.Tag, error)
	GetByID(id int) (models.Tag, error)
	Create(tag models.Tag) (models.Tag, error)
}

type CommentRepository interface {
	Create(comment models.Comment) (models.Comment, error)
	GetByPostID(postID int) ([]models.Comment, error)
}

type LikeRepository interface {
	Toggle(userID int, postID int) (bool, error) // returns true if liked, false if unliked
	GetCountByPostID(postID int) (int, error)
}
