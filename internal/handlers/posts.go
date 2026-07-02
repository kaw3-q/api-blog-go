package handlers

import (
	"encoding/json"
	"github.com/kaw3-q/api-blog-go/internal/auth"
	"github.com/kaw3-q/api-blog-go/internal/middleware"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"github.com/kaw3-q/api-blog-go/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	Repo repository.PostRepository
}

func NewPostHandler(repo repository.PostRepository) *PostHandler {
	return &PostHandler{Repo: repo}
}

func slugify(title string) string {
	title = strings.ToLower(title)
	var sb strings.Builder
	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			sb.WriteRune('-')
		}
	}
	res := sb.String()
	for strings.Contains(res, "--") {
		res = strings.ReplaceAll(res, "--", "-")
	}
	return strings.Trim(res, "-")
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/posts" || r.URL.Path == "/posts/" {
		h.handlePosts(w, r)
		return
	}
	h.handlePostByID(w, r)
}

func (h *PostHandler) handlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		posts, err := h.Repo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	case http.MethodPost:
		var req struct {
			Title           string            `json:"title"`
			Excerpt         *string           `json:"excerpt"`
			Content         string            `json:"content"`
			CoverImage      *string           `json:"cover_image"`
			Status          models.PostStatus `json:"status"`
			Featured        bool              `json:"featured"`
			ReadingTime     *int              `json:"reading_time"`
			MetaTitle       *string           `json:"meta_title"`
			MetaDescription *string           `json:"meta_description"`
			CategoryID      *int              `json:"category_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p := models.Post{
			Title:           req.Title,
			Slug:            slugify(req.Title),
			Excerpt:         req.Excerpt,
			Content:         req.Content,
			CoverImage:      req.CoverImage,
			Status:          req.Status,
			Featured:        req.Featured,
			ReadingTime:     req.ReadingTime,
			MetaTitle:       req.MetaTitle,
			MetaDescription: req.MetaDescription,
			CategoryID:      req.CategoryID,
		}

		if p.Status == "" {
			p.Status = models.PostStatusDraft
		}

		claims, ok := r.Context().Value(middleware.UserKey).(*auth.Claims)
		if ok {
			p.AuthorID = claims.UserID
		}

		newPost, err := h.Repo.Create(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
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
