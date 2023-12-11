package models

import (
	"time"
)

type User struct {
	Id               *int64     `json:"id" db:"id"`
	Firstname        *string    `json:"firstname" db:"firstname"`
	Lastname         *string    `json:"lastname" db:"lastname"`
	Surname          *string    `json:"surname" db:"surname"`
	Email            *string    `json:"email" db:"email" binding:"required"`
	Password         *string    `json:"password" db:"password" binding:"required"`
	TimeZone         *string    `json:"time_zone" db:"time_zone"`
	TimeWorkStart    *time.Time `json:"time_work_start" db:"time_work_start"`
	TimeWorkEnd      *time.Time `json:"time_work_end" db:"time_work_end"`
	ServiceTime      *time.Time `json:"service_time" db:"service_time"`
	DinnerTime       *time.Time `json:"dinner_time" db:"dinner_time"`
	PhoneCountryCode *string    `json:"phone_country_code" db:"phone_country_code"`
	PhoneNumber      *string    `json:"phone_number" db:"phone_number"`
}
