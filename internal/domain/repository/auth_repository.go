package repository

import "github.com/gurodrigues-dev/notifier-app/internal/entity"

type AuthRepository interface {
	CreateToken(token *entity.Token) error
	GetTokenByUser(email string) (*entity.Token, error)
	GetToken(token string) (*entity.Token, error)
	DeleteToken(email string) error
}
