package repository

import (
	"github.com/jmoiron/sqlx"
	"go_api/models"
)

type Authorization interface {
	CreateUser(user models.UserSignUpInput) (int64, error)
	GetUser(username string) (models.User, error)
}

type CompanyDomain interface {
	Create(userId int, list models.CompanyDomain) (int, error)
	GetAll(userId, page, limit int) ([]models.CompanyDomain, error)
	GetById(userId, domainId int) (models.CompanyDomain, error)
	DeleteById(userId, domainId int) error
	Update(userId, domainId int, input models.UpdateCompanyDomainInput) error
}

type User interface {
	GetById(userId int64) (models.User, error)
	Update(userId int64, input models.UserInputUpdate) error
}

type Meeting interface {
	Create(dataMeetingCreate models.MeetingDataCreate) (int64, error)
	GetAll(userId int64, page, limit int) ([]models.MeetingFullData, error)
}

type Repository struct {
	Authorization
	CompanyDomain
	User
	Meeting
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSql(db),
		CompanyDomain: NewCompanyDomainPostgres(db),
		User:          NewUserPostgres(db),
		Meeting:       NewMeetingPostgres(db),
	}
}
