package authentication

import (
	"time"
)

type User struct {
	ID          string     `json:"id" db:"id"`
	Email       string     `json:"email" db:"email"`
	Name        string     `json:"name" db:"name"`
	PhoneNumber string     `json:"phone_number" db:"phone_number"`
	BirthDate   *time.Time `json:"birth_date" db:"birth_date"`
}
