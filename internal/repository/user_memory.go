package repository

import (
	"errors"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"sync"
)

type UserRepository interface {
	GetAll() []models.User
	GetByID(id int) (models.User, error)
	Create(user models.User) models.User
}

type memoryUserRepository struct {
	users  []models.User
	nextID int
	mu     sync.RWMutex
}

func NewMemoryUserRepository() UserRepository {
	return &memoryUserRepository{
		users: []models.User{
			{ID: 1, Username: "admin_master", Email: "admin@blog.com", Role: models.RoleAdmin},
			{ID: 2, Username: "joao_silva", Email: "joao@email.com", Role: models.RoleUser},
		},
		nextID: 3,
	}
}

func (r *memoryUserRepository) GetAll() []models.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.users
}

func (r *memoryUserRepository) GetByID(id int) (models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return models.User{}, errors.New("usuário não encontrado")
}

func (r *memoryUserRepository) Create(user models.User) models.User {
	r.mu.Lock()
	defer r.mu.Unlock()
	user.ID = r.nextID
	r.nextID++
	r.users = append(r.users, user)
	return user
}
