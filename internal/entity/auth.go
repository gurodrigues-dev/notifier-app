package entity

import "github.com/google/uuid"

type Token struct {
	ID        int    `json:"id"`
	Token     string `json:"token"`
	AdminUser string `json:"admin_user" validate:"required,email"`
}

func (t *Token) CreateToken() string {
	t.Token = uuid.NewString()
	return t.Token
}
