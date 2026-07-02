package handlers

import (
	"encoding/json"
	"github.com/kaw3-q/api-blog-go/internal/auth"
	"github.com/kaw3-q/api-blog-go/internal/middleware"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"github.com/kaw3-q/api-blog-go/internal/repository"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	Repo repository.CategoryRepository
}

func NewCategoryHandler(repo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{Repo: repo}
}

func (h *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		categories, err := h.Repo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	case http.MethodPost:
		var req struct {
			Name        string  `json:"name"`
			Slug        string  `json:"slug"`
			Description *string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		slug := req.Slug
		if slug == "" {
			slug = slugify(req.Name)
		}

		cat := models.Category{
			Name:        req.Name,
			Slug:        slug,
			Description: req.Description,
		}

		newCat, err := h.Repo.Create(cat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCat)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

type TagHandler struct {
	Repo repository.TagRepository
}

func NewTagHandler(repo repository.TagRepository) *TagHandler {
	return &TagHandler{Repo: repo}
}

func (h *TagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tags, err := h.Repo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tags)
	case http.MethodPost:
		var req struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		slug := req.Slug
		if slug == "" {
			slug = slugify(req.Name)
		}

		tag := models.Tag{
			Name: req.Name,
			Slug: slug,
		}

		newTag, err := h.Repo.Create(tag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTag)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

type CommentHandler struct {
	Repo repository.CommentRepository
}

func NewCommentHandler(repo repository.CommentRepository) *CommentHandler {
	return &CommentHandler{Repo: repo}
}

func (h *CommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		postIDStr := r.URL.Query().Get("post_id")
		if postIDStr == "" {
			http.Error(w, "post_id é obrigatório para listagem", http.StatusBadRequest)
			return
		}
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "post_id inválido", http.StatusBadRequest)
			return
		}

		comments, err := h.Repo.GetByPostID(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	case http.MethodPost:
		var req struct {
			Content string `json:"content"`
			PostID  int    `json:"post_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := r.Context().Value(middleware.UserKey).(*auth.Claims)
		if !ok {
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		comment := models.Comment{
			Content:  req.Content,
			PostID:   req.PostID,
			AuthorID: claims.UserID,
		}

		newComment, err := h.Repo.Create(comment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newComment)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

type LikeHandler struct {
	Repo repository.LikeRepository
}

func NewLikeHandler(repo repository.LikeRepository) *LikeHandler {
	return &LikeHandler{Repo: repo}
}

func (h *LikeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req struct {
			PostID int `json:"post_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := r.Context().Value(middleware.UserKey).(*auth.Claims)
		if !ok {
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		liked, err := h.Repo.Toggle(claims.UserID, req.PostID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		count, err := h.Repo.GetCountByPostID(req.PostID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"liked": liked,
			"count": count,
		})
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}
