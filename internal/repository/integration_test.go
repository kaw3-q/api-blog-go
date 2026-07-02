package repository

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"github.com/kaw3-q/api-blog-go/prisma/db"
)

func TestDatabaseIntegration(t *testing.T) {
	// Carrega .env a partir da raiz do repositório
	_ = godotenv.Load("../../.env")

	client := db.NewClient()
	if err := client.Connect(); err != nil {
		t.Skip("Skipping integration test: database not available or not configured")
		return
	}
	defer func() {
		_ = client.Disconnect()
	}()

	ctx := context.Background()
	randomSuffix := rand.Intn(100000)

	testUsername := fmt.Sprintf("integration_user_%d", randomSuffix)
	testEmail := fmt.Sprintf("integration_%d@example.com", randomSuffix)
	testCategoryName := fmt.Sprintf("Integration Cat %d", randomSuffix)
	testCategorySlug := fmt.Sprintf("integration-cat-%d", randomSuffix)
	testPostTitle := fmt.Sprintf("Integration Post %d", randomSuffix)
	testPostSlug := fmt.Sprintf("integration-post-%d", randomSuffix)

	// 1. Instancia os repositórios reais do Prisma
	userRepo := NewPrismaUserRepository(client)
	postRepo := NewPrismaPostRepository(client)
	categoryRepo := NewPrismaCategoryRepository(client)
	commentRepo := NewPrismaCommentRepository(client)
	likeRepo := NewPrismaLikeRepository(client)

	// 2. Cria Categoria
	cat, err := categoryRepo.Create(models.Category{
		Name: testCategoryName,
		Slug: testCategorySlug,
	})
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	defer func() {
		_, _ = client.Category.FindUnique(db.Category.ID.Equals(cat.ID)).Delete().Exec(ctx)
	}()

	// 3. Cria Usuário
	usr, err := userRepo.Create(models.User{
		Username: testUsername,
		Email:    testEmail,
		Password: "hashed_integration_password",
		Role:     models.RoleUser,
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer func() {
		_, _ = client.User.FindUnique(db.User.ID.Equals(usr.ID)).Delete().Exec(ctx)
	}()

	// 4. Cria Post
	post, err := postRepo.Create(models.Post{
		Title:      testPostTitle,
		Slug:       testPostSlug,
		Content:    "This is a test post created during integration tests.",
		AuthorID:   usr.ID,
		CategoryID: &cat.ID,
		Status:     models.PostStatusPublished,
	})
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}
	defer func() {
		_, _ = client.Post.FindUnique(db.Post.ID.Equals(post.ID)).Delete().Exec(ctx)
	}()

	// 5. Cria Comentário
	comment, err := commentRepo.Create(models.Comment{
		Content:  "Awesome integration post!",
		AuthorID: usr.ID,
		PostID:   post.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}
	defer func() {
		_, _ = client.Comment.FindUnique(db.Comment.ID.Equals(comment.ID)).Delete().Exec(ctx)
	}()

	// 6. Curte o Post (Toggle Like)
	liked, err := likeRepo.Toggle(usr.ID, post.ID)
	if err != nil {
		t.Fatalf("Failed to toggle like: %v", err)
	}
	if !liked {
		t.Errorf("Expected liked to be true on first toggle")
	}
	defer func() {
		// Garante que limpamos o Like se sobrar no banco
		_, _ = client.Like.FindUnique(db.Like.UserIDPostID(db.Like.UserID.Equals(usr.ID), db.Like.PostID.Equals(post.ID))).Delete().Exec(ctx)
	}()

	// 7. Validação: Comentários e Contagem de Likes
	comments, err := commentRepo.GetByPostID(post.ID)
	if err != nil {
		t.Fatalf("Failed to get comments: %v", err)
	}
	if len(comments) != 1 || comments[0].Content != "Awesome integration post!" {
		t.Errorf("Comments verification failed")
	}

	likeCount, err := likeRepo.GetCountByPostID(post.ID)
	if err != nil {
		t.Fatalf("Failed to get like count: %v", err)
	}
	if likeCount != 1 {
		t.Errorf("Expected 1 like, got %d", likeCount)
	}

	// 8. Descurte (Toggle Like de novo)
	unliked, err := likeRepo.Toggle(usr.ID, post.ID)
	if err != nil {
		t.Fatalf("Failed to toggle unlike: %v", err)
	}
	if unliked {
		t.Errorf("Expected liked to be false on second toggle (unlike)")
	}

	likeCountAfter, err := likeRepo.GetCountByPostID(post.ID)
	if err != nil {
		t.Fatalf("Failed to get like count after unlike: %v", err)
	}
	if likeCountAfter != 0 {
		t.Errorf("Expected 0 likes after unlike, got %d", likeCountAfter)
	}
}
