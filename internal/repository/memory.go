package repository

import (
	"errors"
	"hello-go/internal/models"
	"sync"
)

type PostRepository interface {
	GetAll() []models.Post
	GetByID(id int) (models.Post, error)
	Create(post models.Post) models.Post
}

type memoryPostRepository struct {
	posts  []models.Post
	nextID int
	mu     sync.RWMutex
}

func NewMemoryPostRepository() PostRepository {
	return &memoryPostRepository{
		posts: []models.Post{
			{ID: 1, Title: "Bem-vindo ao meu Blog", Content: "Este é o primeiro post!"},
		},
		nextID: 2,
	}
}

func (r *memoryPostRepository) GetAll() []models.Post {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.posts
}

func (r *memoryPostRepository) GetByID(id int) (models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.posts {
		if p.ID == id {
			return p, nil
		}
	}
	return models.Post{}, errors.New("post não encontrado")
}

func (r *memoryPostRepository) Create(post models.Post) models.Post {
	r.mu.Lock()
	defer r.mu.Unlock()
	post.ID = r.nextID
	r.nextID++
	r.posts = append(r.posts, post)
	return post
}
