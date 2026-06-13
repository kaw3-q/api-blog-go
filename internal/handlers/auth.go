package handlers

import (
	"encoding/json"
	"github.com/kaw3-q/api-blog-go/internal/auth"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"github.com/kaw3-q/api-blog-go/internal/repository"
	"net/http"
)

type AuthHandler struct {
	UserRepo repository.UserRepository
}

func NewAuthHandler(repo repository.UserRepository) *AuthHandler {
	return &AuthHandler{UserRepo: repo}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashed, _ := auth.HashPassword(u.Password)
	u.Password = hashed
	if u.Role == "" {
		u.Role = models.RoleUser
	}

	newUser := h.UserRepo.Create(u)
	json.NewEncoder(w).Encode(newUser)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.GetByEmail(req.Email)
	if err != nil || !auth.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateToken(user)
	json.NewEncoder(w).Encode(models.LoginResponse{
		Token: token,
		User:  user,
	})
}
