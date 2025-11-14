package auth_domain

import (
	"fmt"
	"time"

	validate "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID             uuid.UUID `json:"user_id"`
	Email              string    `json:"email"`
	Password           string    `json:"password,omitempty"`
	EncryptedPassword  string    `json:"-"`
	RefreshToken       string    `json:"-"`
	RefreshTokenExpire time.Time `json:"-"`
}

func NewUser(email string, password string) (*User, error) {
	u := &User{
		Email: email,
		Password: password,
	}

	if err := u.ValidateUser(); err != nil {
		return nil, fmt.Errorf("invalid user: %w", err)
	}

	if err := u.GetEncryptPassword(u.Password); err != nil {
		return nil, err
	}

	u.Password = ""

	return u, nil
}

func (u *User) ValidateUser() error {
	return validate.ValidateStruct(
		u,
		validate.Field(&u.Email, validate.Required, is.Email),
		validate.Field(&u.Password, validate.By(requiredIf(u.EncryptedPassword == "")), validate.Length(8, 100)),
	)
}

func (u *User) GetEncryptPassword(password string) error {
	if len(password) > 0 {
		enc, err := encryptPassword(password)
		if err != nil {
			return fmt.Errorf("failed to get ecnrypted password: %w", err)
		}

		u.EncryptedPassword = enc
	}
	return nil
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt password: %w", err)
	}

	return string(b), nil
}

func requiredIf(cond bool) validate.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validate.Validate(value, validate.Required)
		}

		return nil
	}
}
