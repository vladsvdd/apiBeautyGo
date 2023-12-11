package service

import (
	"go_api/models"
	"go_api/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.UserSignUpInput) (int64, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int64, error)
}

type User interface {
	GetById(userId int64) (models.User, error)
	Update(userId int64, input models.UserInputUpdate) error
}

type Service struct {
	Authorization
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User),
	}
}
