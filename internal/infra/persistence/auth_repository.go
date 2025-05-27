package persistence

import (
	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
)

type AuthRepositoryImpl struct {
	Postgres contracts.PostgresIface
}

func (ar AuthRepositoryImpl) CreateToken(token *entity.Token) error {
	return ar.Postgres.Client().Create(token).Error
}

func (ar AuthRepositoryImpl) GetToken(token string) (*entity.Token, error) {
	var validToken entity.Token
	err := ar.Postgres.Client().Where("token = ?", token).First(&validToken).Error
	if err != nil {
		return nil, err
	}
	return &validToken, nil
}

func (ar AuthRepositoryImpl) GetTokenByUser(email string) (*entity.Token, error) {
	var token entity.Token
	err := ar.Postgres.Client().Where("admin_user = ?", email).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (ar AuthRepositoryImpl) DeleteToken(email string) error {
	return ar.Postgres.Client().Delete(&entity.Token{}, "admin_user = ?", email).Error
}
