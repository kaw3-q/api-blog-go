package models

import "time"

type Role string

const (
	RoleUser   Role = "USER"
	RoleEditor Role = "EDITOR"
	RoleAdmin  Role = "ADMIN"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Oculto no JSON
	Role      Role      `json:"role"`
	Avatar    *string   `json:"avatar,omitempty"`
	Bio       *string   `json:"bio,omitempty"`
	Github    *string   `json:"github,omitempty"`
	Linkedin  *string   `json:"linkedin,omitempty"`
	Website   *string   `json:"website,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
