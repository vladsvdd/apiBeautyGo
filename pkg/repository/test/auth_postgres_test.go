package test

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"go_api/models"
	"go_api/pkg/repository"
	"testing"
	"time"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewAuthSql(db)

	testCasesData := []models.UserSignUpInput{
		{
			Firstname:        "Ivan",
			Lastname:         "Borshev",
			Surname:          "4",
			Email:            "ivan@mail.ru",
			Password:         "...",
			PhoneCountryCode: "7",
			PhoneNumber:      "9845856574",
		},
	}

	tests := []struct {
		name    string
		mock    func()
		input   models.UserSignUpInput
		want    int64
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id FROM users WHERE phone_country_code=\\$1 AND phone_number=\\$2 OR email=\\$3 LIMIT 1").
					WithArgs(testCasesData[0].PhoneCountryCode, testCasesData[0].PhoneNumber, testCasesData[0].Email).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("INSERT INTO users (.+) RETURNING id").
					WithArgs(testCasesData[0].Firstname, testCasesData[0].Lastname, testCasesData[0].Surname, testCasesData[0].Email, testCasesData[0].Password, testCasesData[0].PhoneCountryCode, testCasesData[0].PhoneNumber).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			input: testCasesData[0],
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateUser(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("expectations were not met: %s", err)
			}
		})
	}
}

func parseDate(date string) time.Time {
	// Парсим строку в time.Time
	parsedTime, _ := time.Parse(time.RFC3339, date)
	return parsedTime
}

func TestAuthPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewAuthSql(db)

	testCasesData := []models.User{
		{
			Id:            5,
			Firstname:     "1my name",
			Lastname:      "1фамилия",
			Surname:       "отчество",
			Email:         "t.t@mail",
			Password:      "$2a$10$g1AOohyc3KD3.38peLCQRmhy2LNCwCZP...iQ0grhugyi",
			TimeZone:      "+1",
			TimeWorkStart: parseDate("0001-01-01T09:00:00Z"),
			TimeWorkEnd:   parseDate("0001-01-01T18:00:00Z"),
			ServiceTime:   parseDate("0001-01-01T01:00:00Z"),
			DinnerTime:    parseDate("0001-01-01T12:00:00Z"),
		},
		{
			Id:            0,
			Firstname:     "1my name",
			Lastname:      "1фамилия",
			Surname:       "отчество",
			Email:         "test3@mail.ru",
			Password:      "$2a$10$g1AOohyc3KD3B4cmWZ38peLCQ...yiQ0grhugyi",
			TimeZone:      "+1",
			TimeWorkStart: parseDate("0001-01-01T09:00:00Z"),
			TimeWorkEnd:   parseDate("0001-01-01T18:00:00Z"),
			ServiceTime:   parseDate("0001-01-01T01:00:00Z"),
			DinnerTime:    parseDate("0001-01-01T12:00:00Z"),
		},
	}

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    models.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectQuery("SELECT id, password_hash AS password FROM users WHERE email=\\$1 LIMIT 1").
					WithArgs(testCasesData[0].Email).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"firstname",
						"lastname",
						"surname",
						"email",
						"password",
						"time_zone",
						"time_work_start",
						"time_work_end",
						"service_time",
						"dinner_time",
					}).AddRow(
						testCasesData[0].Id,
						testCasesData[0].Firstname,
						testCasesData[0].Lastname,
						testCasesData[0].Surname,
						testCasesData[0].Email,
						testCasesData[0].Password,
						testCasesData[0].TimeZone,
						testCasesData[0].TimeWorkStart,
						testCasesData[0].TimeWorkEnd,
						testCasesData[0].ServiceTime,
						testCasesData[0].DinnerTime,
					))

			},
			input: testCasesData[0].Email,
			want:  testCasesData[0],
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery("SELECT id, password_hash AS password FROM users WHERE email=\\$1 LIMIT 1").
					WithArgs("t@mail.ru").
					WillReturnError(sql.ErrNoRows)
			},
			input:   "t@mail.ru",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUser(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
