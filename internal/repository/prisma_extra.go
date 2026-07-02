package repository

import (
	"context"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"github.com/kaw3-q/api-blog-go/prisma/db"
)

// CategoryRepository implementation
type prismaCategoryRepository struct {
	client *db.PrismaClient
}

func NewPrismaCategoryRepository(client *db.PrismaClient) CategoryRepository {
	return &prismaCategoryRepository{client: client}
}

func mapCategory(c *db.CategoryModel) models.Category {
	var desc *string
	if val, ok := c.Description(); ok {
		desc = &val
	}
	return models.Category{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: desc,
	}
}

func (r *prismaCategoryRepository) GetAll() ([]models.Category, error) {
	categories, err := r.client.Category.FindMany().Exec(context.Background())
	if err != nil {
		return nil, err
	}
	var res []models.Category
	for _, c := range categories {
		res = append(res, mapCategory(&c))
	}
	return res, nil
}

func (r *prismaCategoryRepository) GetByID(id int) (models.Category, error) {
	c, err := r.client.Category.FindUnique(
		db.Category.ID.Equals(id),
	).Exec(context.Background())
	if err != nil {
		return models.Category{}, err
	}
	return mapCategory(c), nil
}

func (r *prismaCategoryRepository) Create(category models.Category) (models.Category, error) {
	opts := []db.CategorySetParam{}
	if category.Description != nil {
		opts = append(opts, db.Category.Description.Set(*category.Description))
	}
	c, err := r.client.Category.CreateOne(
		db.Category.Name.Set(category.Name),
		db.Category.Slug.Set(category.Slug),
		opts...,
	).Exec(context.Background())
	if err != nil {
		return models.Category{}, err
	}
	return mapCategory(c), nil
}

// TagRepository implementation
type prismaTagRepository struct {
	client *db.PrismaClient
}

func NewPrismaTagRepository(client *db.PrismaClient) TagRepository {
	return &prismaTagRepository{client: client}
}

func mapTag(t *db.TagModel) models.Tag {
	return models.Tag{
		ID:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}
}

func (r *prismaTagRepository) GetAll() ([]models.Tag, error) {
	tags, err := r.client.Tag.FindMany().Exec(context.Background())
	if err != nil {
		return nil, err
	}
	var res []models.Tag
	for _, t := range tags {
		res = append(res, mapTag(&t))
	}
	return res, nil
}

func (r *prismaTagRepository) GetByID(id int) (models.Tag, error) {
	t, err := r.client.Tag.FindUnique(
		db.Tag.ID.Equals(id),
	).Exec(context.Background())
	if err != nil {
		return models.Tag{}, err
	}
	return mapTag(t), nil
}

func (r *prismaTagRepository) Create(tag models.Tag) (models.Tag, error) {
	t, err := r.client.Tag.CreateOne(
		db.Tag.Name.Set(tag.Name),
		db.Tag.Slug.Set(tag.Slug),
	).Exec(context.Background())
	if err != nil {
		return models.Tag{}, err
	}
	return mapTag(t), nil
}

// CommentRepository implementation
type prismaCommentRepository struct {
	client *db.PrismaClient
}

func NewPrismaCommentRepository(client *db.PrismaClient) CommentRepository {
	return &prismaCommentRepository{client: client}
}

func mapComment(c *db.CommentModel) models.Comment {
	return models.Comment{
		ID:        c.ID,
		Content:   c.Content,
		AuthorID:  c.AuthorID,
		PostID:    c.PostID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (r *prismaCommentRepository) Create(comment models.Comment) (models.Comment, error) {
	c, err := r.client.Comment.CreateOne(
		db.Comment.Content.Set(comment.Content),
		db.Comment.Author.Link(db.User.ID.Equals(comment.AuthorID)),
		db.Comment.Post.Link(db.Post.ID.Equals(comment.PostID)),
	).Exec(context.Background())
	if err != nil {
		return models.Comment{}, err
	}
	return mapComment(c), nil
}

func (r *prismaCommentRepository) GetByPostID(postID int) ([]models.Comment, error) {
	comments, err := r.client.Comment.FindMany(
		db.Comment.PostID.Equals(postID),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	var res []models.Comment
	for _, c := range comments {
		res = append(res, mapComment(&c))
	}
	return res, nil
}

// LikeRepository implementation
type prismaLikeRepository struct {
	client *db.PrismaClient
}

func NewPrismaLikeRepository(client *db.PrismaClient) LikeRepository {
	return &prismaLikeRepository{client: client}
}

func (r *prismaLikeRepository) Toggle(userID int, postID int) (bool, error) {
	ctx := context.Background()
	// Check if already liked
	like, err := r.client.Like.FindUnique(
		db.Like.UserIDPostID(
			db.Like.UserID.Equals(userID),
			db.Like.PostID.Equals(postID),
		),
	).Exec(ctx)

	if err == nil && like != nil {
		// Unlike
		_, err = r.client.Like.FindUnique(
			db.Like.UserIDPostID(
				db.Like.UserID.Equals(userID),
				db.Like.PostID.Equals(postID),
			),
		).Delete().Exec(ctx)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// Like
	_, err = r.client.Like.CreateOne(
		db.Like.User.Link(db.User.ID.Equals(userID)),
		db.Like.Post.Link(db.Post.ID.Equals(postID)),
	).Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *prismaLikeRepository) GetCountByPostID(postID int) (int, error) {
	count, err := r.client.Like.FindMany(
		db.Like.PostID.Equals(postID),
	).Exec(context.Background())
	if err != nil {
		return 0, err
	}
	return len(count), nil
}
