package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go_api/models"
)

type AuthSql struct {
	db *sqlx.DB
}

func NewAuthSql(db *sqlx.DB) *AuthSql {
	return &AuthSql{db: db}
}

// CreateUser Регистрация нового пользователя
func (r *AuthSql) CreateUser(user models.UserSignUpInput) (int64, error) {
	// Валидация входных данных
	if err := validateUserSignUpInput(user); err != nil {
		return 0, err
	}

	transaction, err := r.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return 0, err
	}

	// Проверка наличия телефона или email в БД
	if hasPhoneOrEmail(transaction, user.PhoneCountryCode, user.PhoneNumber, user.Email) {
		transaction.Rollback()
		return 0, errors.New(fmt.Sprint("Пользователь с таким телефоном или email уже зарегистрирован"))
	}

	//добавляем пользователя
	var id int64
	query := fmt.Sprintf(`
		INSERT INTO
		    %s (firstname,
				lastname,
				surname,
				email,
				password_hash,
		        phone_country_code,
		        phone_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`, usersTable)

	row := transaction.QueryRow(query,
		user.Firstname,
		user.Lastname,
		user.Surname,
		user.Email,
		user.Password,
		user.PhoneCountryCode,
		user.PhoneNumber)
	if err := row.Scan(&id); err != nil {
		transaction.Rollback()
		logrus.Error(err)
		return 0, errors.New(fmt.Sprint("Пользователь с таким email или телефоном уже зарегистрирован"))
	}

	return id, transaction.Commit()
}

// Валидация входных данных пользователя
func validateUserSignUpInput(user models.UserSignUpInput) error {
	if user.Email == "" {
		return errors.New(fmt.Sprint("Заполните email"))
	}

	if user.PhoneCountryCode == "" {
		return errors.New(fmt.Sprint("Заполните телефонный код страны"))
	}

	if user.PhoneNumber == "" {
		return errors.New(fmt.Sprint("Заполните номер телефона"))
	}

	return nil
}

func hasPhoneOrEmail(tx *sql.Tx, phoneCountryCode, phoneNumber, email string) bool {
	query := fmt.Sprintf(`
		SELECT id 
		FROM %s
		WHERE
		    phone_country_code=$1 AND phone_number=$2
		OR
		    email=$3
		LIMIT 1
	`, usersTable)

	var id int64
	row := tx.QueryRow(query, phoneCountryCode, phoneNumber, email)
	if err := row.Scan(&id); err == sql.ErrNoRows {
		return false
	}
	return true
}

func (r *AuthSql) GetUser(email string) (user models.User, err error) {
	if email == "" {
		return user, errors.New(fmt.Sprint("Заполните email"))
	}

	query := fmt.Sprintf("SELECT id, password_hash AS password FROM %s WHERE email=$1 LIMIT 1", usersTable)
	err = r.db.Get(&user, query, email)

	return user, err
}
