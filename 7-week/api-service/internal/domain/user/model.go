package user

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleUser
}

type User struct {
	ID uuid.UUID

	FirstName string `validate:"required,min=2,max=50"`
	LastName  string `validate:"required,min=2,max=50"`
	Email     string `validate:"required,email,max=254"`
	Role      Role   `validate:"required"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(firstName, lastName, email string, role Role) (*User, error) {
	u := &User{
		ID:   uuid.New(),
		Role: role,
	}
	if err := u.UpdateProfile(firstName, lastName, email); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) UpdateProfile(firstName, lastName, email string) error {
	u.FirstName = strings.TrimSpace(firstName)
	u.LastName = strings.TrimSpace(lastName)
	u.Email = NormalizeEmail(email)

	return u.Validate()
}

func (u *User) Validate() error {
	if err := mapValidationError(validate.Struct(u)); err != nil {
		return err
	}
	if !u.Role.IsValid() {
		return ErrInvalidRole
	}

	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func mapValidationError(err error) error {
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	for _, fieldErr := range validationErrors {
		if mapped, ok := fieldErrors[fieldErr.Field()]; ok {
			return mapped
		}
	}

	return err
}
