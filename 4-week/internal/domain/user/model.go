package user

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

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
	Email     string `validate:"required,email"`
	Role      Role   `validate:"required"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(id uuid.UUID, firstName string, lastName string, email string, role Role) (*User, error) {
	u := &User{
		ID:        id,
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Role:      role,
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
	if err := mapValidationError(validate.Struct(u)); err != nil {
		return err
	}

	if !u.Role.IsValid() {
		return ErrInvalidRole
	}

	return nil
}

func ValidatePassword(password string) error {
	if err := validate.Var(password, "required"); err != nil {
		return ErrInvalidPassword
	}

	return nil
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
		if err, ok := fieldErrors[fieldErr.Field()]; ok {
			return err
		}
	}

	return err
}
