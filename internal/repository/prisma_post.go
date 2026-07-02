
package repository

import (
	"context"
	"time"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"github.com/kaw3-q/api-blog-go/prisma/db"
)

type prismaPostRepository struct {
	client *db.PrismaClient
}

func NewPrismaPostRepository(client *db.PrismaClient) PostRepository {
	return &prismaPostRepository{client: client}
}

func mapPost(p *db.PostModel) models.Post {
	var excerpt, coverImage *string
	var categoryID *int
	var readingTime *int
	var metaTitle, metaDescription *string
	var publishedAt *time.Time

	if val, ok := p.Excerpt(); ok {
		excerpt = &val
	}
	if val, ok := p.CoverImage(); ok {
		coverImage = &val
	}
	if val, ok := p.CategoryID(); ok {
		categoryID = &val
	}
	if val, ok := p.ReadingTime(); ok {
		readingTime = &val
	}
	if val, ok := p.MetaTitle(); ok {
		metaTitle = &val
	}
	if val, ok := p.MetaDescription(); ok {
		metaDescription = &val
	}
	if val, ok := p.PublishedAt(); ok {
		publishedAt = &val
	}

	return models.Post{
		ID:              p.ID,
		Title:           p.Title,
		Slug:            p.Slug,
		Excerpt:         excerpt,
		Content:         p.Content,
		CoverImage:      coverImage,
		Status:          models.PostStatus(p.Status),
		Featured:        p.Featured,
		Views:           p.Views,
		ReadingTime:     readingTime,
		MetaTitle:       metaTitle,
		MetaDescription: metaDescription,
		AuthorID:        p.AuthorID,
		CategoryID:      categoryID,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
		PublishedAt:     publishedAt,
	}
}

func (r *prismaPostRepository) GetAll() ([]models.Post, error) {
	posts, err := r.client.Post.FindMany().Exec(context.Background())
	if err != nil {
		return nil, err
	}
	var res []models.Post
	for _, p := range posts {
		res = append(res, mapPost(&p))
	}
	return res, nil
}

func (r *prismaPostRepository) GetByID(id int) (models.Post, error) {
	p, err := r.client.Post.FindUnique(
		db.Post.ID.Equals(id),
	).Exec(context.Background())
	if err != nil {
		return models.Post{}, err
	}
	return mapPost(p), nil
}

func (r *prismaPostRepository) Create(post models.Post) (models.Post, error) {
	opts := []db.PostSetParam{
		db.Post.Status.Set(db.PostStatus(post.Status)),
		db.Post.Featured.Set(post.Featured),
		db.Post.Views.Set(post.Views),
	}
	if post.Excerpt != nil {
		opts = append(opts, db.Post.Excerpt.Set(*post.Excerpt))
	}
	if post.CoverImage != nil {
		opts = append(opts, db.Post.CoverImage.Set(*post.CoverImage))
	}
	if post.ReadingTime != nil {
		opts = append(opts, db.Post.ReadingTime.Set(*post.ReadingTime))
	}
	if post.MetaTitle != nil {
		opts = append(opts, db.Post.MetaTitle.Set(*post.MetaTitle))
	}
	if post.MetaDescription != nil {
		opts = append(opts, db.Post.MetaDescription.Set(*post.MetaDescription))
	}
	if post.CategoryID != nil {
		opts = append(opts, db.Post.CategoryID.Set(*post.CategoryID))
	}
	if post.PublishedAt != nil {
		opts = append(opts, db.Post.PublishedAt.Set(*post.PublishedAt))
	}

	p, err := r.client.Post.CreateOne(
		db.Post.Title.Set(post.Title),
		db.Post.Slug.Set(post.Slug),
		db.Post.Content.Set(post.Content),
		db.Post.Author.Link(db.User.ID.Equals(post.AuthorID)),
		opts...,
	).Exec(context.Background())
	if err != nil {
		return models.Post{}, err
	}
	return mapPost(p), nil
}
