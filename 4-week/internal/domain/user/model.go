package user

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type User struct {
	ID uuid.UUID

	FirstName string `validate:"required,min=2,max=50"`
	LastName  string `validate:"required,min=2,max=50"`
	Email     string `validate:"required,email"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(id uuid.UUID, firstName string, lastName string, email string) (*User, error) {
	u := &User{
		ID:        id,
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Email:     strings.TrimSpace(email),
	}

	if u.ID == uuid.Nil {
		return nil, ErrInvalidID
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u User) Validate() error {
	return mapValidationError(validate.Struct(u))
}

func mapValidationError(err error) error {
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	for _, fieldErr := range validationErrors {
		switch fieldErr.Field() {
		case "FirstName":
			return ErrInvalidFirstName
		case "LastName":
			return ErrInvalidLastName
		case "Email":
			return ErrInvalidEmail
		}
	}

	return err
}
