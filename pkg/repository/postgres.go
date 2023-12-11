package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx" // Импорт пакета sqlx для работы с базами данных SQL.
)

const (
	usersTable                     = "users"
	phoneNumbersTable              = "phone_numbers"
	phoneNumbersTypeTable          = "phone_number_type"
	companyDomainsTable            = "company_domains"
	linkUsersToCompanyDomainsTable = "link_users_to_company_domains"
	todoListsTable                 = "todo_lists"
	usersListsTable                = "users_lists"
	todoItemsTable                 = "todo_items"
	listsItemsTable                = "lists_items"
	clientTable                    = "clients"
	meetingTable                   = "meetings"
	meetingStatusesTable           = "meeting_statuses"
)

// Config представляет конфигурацию подключения к базе данных PostgreSQL.
type Config struct {
	Host     string // Хост базы данных.
	Port     string // Порт базы данных.
	Username string // Имя пользователя базы данных.
	Password string // Пароль для подключения к базе данных.
	DBName   string // Имя базы данных.
	SSLMode  string // Режим SSL для подключения к базе данных.
}

// NewPostgresDB создает новое подключение к базе данных PostgreSQL на основе переданной конфигурации.
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)) // Открытие подключения к базе данных PostgreSQL.
	if err != nil {
		return nil, err
	}

	err = db.Ping() // Проверка доступности базы данных.
	if err != nil {
		return nil, err
	}

	return db, nil // Возвращение экземпляра подключения к базе данных.
}
