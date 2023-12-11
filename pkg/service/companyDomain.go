package service

import (
	go_api "go_api/models"
	"go_api/pkg/repository"
)

type CompanyDomainService struct {
	repo repository.CompanyDomain
}

func NewCompanyDomainService(repo repository.CompanyDomain) *CompanyDomainService {
	return &CompanyDomainService{repo: repo}
}

func (s *CompanyDomainService) Create(userId int, domain go_api.CompanyDomain) (int, error) {
	return s.repo.Create(userId, domain)
}

func (s *CompanyDomainService) GetAll(userId, page, limit int) ([]go_api.CompanyDomain, error) {
	return s.repo.GetAll(userId, page, limit)
}

func (s *CompanyDomainService) GetById(userId, domainId int) (go_api.CompanyDomain, error) {
	return s.repo.GetById(userId, domainId)
}

func (s *CompanyDomainService) DeleteById(userId, domainId int) error {
	return s.repo.DeleteById(userId, domainId)
}

func (s *CompanyDomainService) Update(userId, domainId int, input go_api.UpdateCompanyDomainInput) error {
	err := input.Validate()
	if err != nil {
		return err
	}

	return s.repo.Update(userId, domainId, input)
}
