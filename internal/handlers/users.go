package handlers

import (
	"encoding/json"
	"github.com/kaw3-q/api-blog-go/internal/auth"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"github.com/kaw3-q/api-blog-go/internal/repository"
	"net/http"
	"strconv"
)

type UserHandler struct {
	Repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/users" || r.URL.Path == "/users/" {
		h.handleUsers(w, r)
		return
	}
	h.handleUserByID(w, r)
}

func (h *UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, err := h.Repo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	case http.MethodPost:
		var req struct {
			Username string      `json:"username"`
			Email    string      `json:"email"`
			Password string      `json:"password"`
			Role     models.Role `json:"role"`
			Avatar   *string     `json:"avatar"`
			Bio      *string     `json:"bio"`
			Github   *string     `json:"github"`
			Linkedin *string     `json:"linkedin"`
			Website  *string     `json:"website"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		role := req.Role
		if role == "" {
			role = models.RoleUser
		}

		hashed, _ := auth.HashPassword(req.Password)

		u := models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: hashed,
			Role:     role,
			Avatar:   req.Avatar,
			Bio:      req.Bio,
			Github:   req.Github,
			Linkedin: req.Linkedin,
			Website:  req.Website,
		}

		newUser, err := h.Repo.Create(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
