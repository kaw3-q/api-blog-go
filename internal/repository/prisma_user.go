package repository

import (
	"context"
	"github.com/kaw3-q/api-blog-go/internal/models"
	"github.com/kaw3-q/api-blog-go/prisma/db"
)

type prismaUserRepository struct {
	client *db.PrismaClient
}

func NewPrismaUserRepository(client *db.PrismaClient) UserRepository {
	return &prismaUserRepository{client: client}
}

func mapUser(u *db.UserModel) models.User {
	var avatar, bio, github, linkedin, website *string
	if val, ok := u.Avatar(); ok {
		avatar = &val
	}
	if val, ok := u.Bio(); ok {
		bio = &val
	}
	if val, ok := u.Github(); ok {
		github = &val
	}
	if val, ok := u.Linkedin(); ok {
		linkedin = &val
	}
	if val, ok := u.Website(); ok {
		website = &val
	}
	return models.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Role:      models.Role(u.Role),
		Avatar:    avatar,
		Bio:       bio,
		Github:    github,
		Linkedin:  linkedin,
		Website:   website,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (r *prismaUserRepository) GetByEmail(email string) (models.User, error) {
	u, err := r.client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(context.Background())
	if err != nil {
		return models.User{}, err
	}
	return mapUser(u), nil
}

func (r *prismaUserRepository) GetByID(id int) (models.User, error) {
	u, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(context.Background())
	if err != nil {
		return models.User{}, err
	}
	return mapUser(u), nil
}

func (r *prismaUserRepository) Create(user models.User) (models.User, error) {
	opts := []db.UserSetParam{
		db.User.Role.Set(db.Role(user.Role)),
	}
	if user.Avatar != nil {
		opts = append(opts, db.User.Avatar.Set(*user.Avatar))
	}
	if user.Bio != nil {
		opts = append(opts, db.User.Bio.Set(*user.Bio))
	}
	if user.Github != nil {
		opts = append(opts, db.User.Github.Set(*user.Github))
	}
	if user.Linkedin != nil {
		opts = append(opts, db.User.Linkedin.Set(*user.Linkedin))
	}
	if user.Website != nil {
		opts = append(opts, db.User.Website.Set(*user.Website))
	}

	u, err := r.client.User.CreateOne(
		db.User.Username.Set(user.Username),
		db.User.Email.Set(user.Email),
		db.User.Password.Set(user.Password),
		opts...,
	).Exec(context.Background())
	if err != nil {
		return models.User{}, err
	}
	return mapUser(u), nil
}

func (r *prismaUserRepository) GetAll() ([]models.User, error) {
	users, err := r.client.User.FindMany().Exec(context.Background())
	if err != nil {
		return nil, err
	}
	var res []models.User
	for _, u := range users {
		res = append(res, mapUser(&u))
	}
	return res, nil
}
