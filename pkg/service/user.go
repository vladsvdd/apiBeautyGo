package service

import (
	"go_api/models"
	"go_api/pkg/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

//func (s *CompanyDomainService) Create(userId int, domain go_api.CompanyDomain) (int, error) {
//	return s.repo.Create(userId, domain)
//}
//
//func (s *CompanyDomainService) GetAll(userId, page, limit int) ([]go_api.CompanyDomain, error) {
//	return s.repo.GetAll(userId, page, limit)
//}

func (s *UserService) GetById(userId int64) (models.User, error) {
	return s.repo.GetById(userId)
}

//func (s *CompanyDomainService) DeleteById(userId, domainId int) error {
//	return s.repo.DeleteById(userId, domainId)
//}

func (s *UserService) Update(userId int64, input models.UserInputUpdate) error {
	return s.repo.Update(userId, input)
}
