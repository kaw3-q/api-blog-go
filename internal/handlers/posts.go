package handlers

import (
	"encoding/json"
	"hello-go/internal/models"
	"hello-go/internal/repository"
	"net/http"
	"strconv"
)

type PostHandler struct {
	Repo repository.PostRepository
}

func NewPostHandler(repo repository.PostRepository) *PostHandler {
	return &PostHandler{Repo: repo}
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Roteamento simples dentro do handler para manter o exemplo conciso
	if r.URL.Path == "/posts" || r.URL.Path == "/posts/" {
		h.handlePosts(w, r)
		return
	}

	// Rota para post individual: /posts/{id}
	h.handlePostByID(w, r)
}

func (h *PostHandler) handlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		posts := h.Repo.GetAll()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	case http.MethodPost:
		var p models.Post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newPost := h.Repo.Create(p)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPost)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (h *PostHandler) handlePostByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/posts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	post, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
