package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go_api/models"
	"strings"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetById(userId int64) (models.User, error) {
	var user models.User

	query := fmt.Sprintf(`
		SELECT
		    u.firstname,
			u.lastname,
			u.surname,
			u.email,
			u.time_zone,
			u.time_work_start,
			u.time_work_end,
			u.service_time,
			u.dinner_time
		FROM
		    %s AS u
		WHERE
		    u.id=$1
		LIMIT 1
	   `, usersTable)
	err := r.db.Get(&user, query, userId)

	return user, err
}

func (r *UserPostgres) Update(userId int64, input models.UserInputUpdate) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Firstname != nil {
		setValues = append(setValues, fmt.Sprintf("firstname=$%d", argId))
		args = append(args, *input.Firstname)
		argId++
	}

	if input.Lastname != nil {
		setValues = append(setValues, fmt.Sprintf("lastname=$%d", argId))
		args = append(args, *input.Lastname)
		argId++
	}

	if input.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
		args = append(args, *input.Surname)
		argId++
	}

	if input.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, *input.Email)
		argId++
	}

	if input.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password=$%d", argId))
		args = append(args, *input.Password)
		argId++
	}

	if input.TimeZone != nil {
		setValues = append(setValues, fmt.Sprintf("time_zone=$%d", argId))
		args = append(args, *input.TimeZone)
		argId++
	}

	if input.TimeWorkStart != nil {
		setValues = append(setValues, fmt.Sprintf("time_work_start=$%d", argId))
		args = append(args, *input.TimeWorkStart)
		argId++
	}

	if input.TimeWorkEnd != nil {
		setValues = append(setValues, fmt.Sprintf("time_work_end=$%d", argId))
		args = append(args, *input.TimeWorkEnd)
		argId++
	}

	if input.ServiceTime != nil {
		setValues = append(setValues, fmt.Sprintf("service_time=$%d", argId))
		args = append(args, *input.ServiceTime)
		argId++
	}

	if input.DinnerTime != nil {
		setValues = append(setValues, fmt.Sprintf("dinner_time=$%d", argId))
		args = append(args, *input.DinnerTime)
		argId++
	}

	// title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`
			UPDATE %s
			SET %s 
			WHERE
			    id = $%d
	`, usersTable, setQuery, argId)

	args = append(args, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return err
}

func hasEmail(db *sqlx.DB, email string) (has bool) {
	var id int64
	query := fmt.Sprintf(`
		SELECT id
		FROM %s
		WHERE
		    email=$1
		LIMIT 1
	`, usersTable)

	row := db.QueryRow(query, email)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			// Запись не найдена, всё в порядке
			return false
		}
	}

	if id == 0 {
		return false
	}

	// Запись найдена
	return true
}
