package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaw3-q/api-blog-go/internal/auth"
	"github.com/kaw3-q/api-blog-go/internal/middleware"
	"github.com/kaw3-q/api-blog-go/internal/models"
)

// MockPostRepository implements repository.PostRepository
type mockPostRepository struct {
	posts []models.Post
}

func (m *mockPostRepository) GetAll() ([]models.Post, error) {
	return m.posts, nil
}

func (m *mockPostRepository) GetByID(id int) (models.Post, error) {
	for _, p := range m.posts {
		if p.ID == id {
			return p, nil
		}
	}
	return models.Post{}, http.ErrNoLocation
}

func (m *mockPostRepository) Create(post models.Post) (models.Post, error) {
	post.ID = len(m.posts) + 1
	m.posts = append(m.posts, post)
	return post, nil
}

func TestPostHandler_GetAll(t *testing.T) {
	mockRepo := &mockPostRepository{
		posts: []models.Post{
			{ID: 1, Title: "Post One", Slug: "post-one", Content: "Content One"},
			{ID: 2, Title: "Post Two", Slug: "post-two", Content: "Content Two"},
		},
	}
	handler := NewPostHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	var posts []models.Post
	if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}
}

func TestPostHandler_Create(t *testing.T) {
	mockRepo := &mockPostRepository{}
	handler := NewPostHandler(mockRepo)

	input := map[string]interface{}{
		"title":   "Hello World",
		"content": "This is a test post.",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Authenticate the request by putting claims in context
	claims := &auth.Claims{UserID: 42}
	ctx := context.WithValue(req.Context(), middleware.UserKey, claims)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", rr.Code)
	}

	var newPost models.Post
	if err := json.NewDecoder(rr.Body).Decode(&newPost); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if newPost.Title != "Hello World" {
		t.Errorf("Expected title 'Hello World', got '%s'", newPost.Title)
	}

	if newPost.Slug != "hello-world" {
		t.Errorf("Expected slug 'hello-world', got '%s'", newPost.Slug)
	}

	if newPost.AuthorID != 42 {
		t.Errorf("Expected AuthorID 42, got %d", newPost.AuthorID)
	}
}
