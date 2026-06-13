package repository

import (
	"github.com/kaw3-q/api-blog-go/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(email string) (models.User, error)
	GetByID(id uint) (models.User, error)
	Create(user models.User) models.User
	GetAll() []models.User
}

type sqlUserRepository struct {
	db *gorm.DB
}

func NewSQLUserRepository(db *gorm.DB) UserRepository {
	return &sqlUserRepository{db: db}
}

func (r *sqlUserRepository) GetByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *sqlUserRepository) GetByID(id uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *sqlUserRepository) Create(user models.User) models.User {
	r.db.Create(&user)
	return user
}

func (r *sqlUserRepository) GetAll() []models.User {
	var users []models.User
	r.db.Find(&users)
	return users
}
