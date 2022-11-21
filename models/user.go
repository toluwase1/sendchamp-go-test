package models

import (
	"errors"
	"fmt"
	goval "github.com/go-passwd/validator"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/leebenson/conform"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Name           string `json:"name" binding:"required,min=2"`
	Email          string `json:"email" gorm:"unique;not null" binding:"required,email"`
	PhoneNumber    string `json:"phone_number" gorm:"unique;default:null" binding:"required,e164"`
	Password       string `json:"password" gorm:"-" binding:"required,min=8,max=15"`
	HashedPassword string `json:"-" gorm:"password"`
	IsEmailActive  bool   `json:"-"`
}

type Model struct {
	ID        uint  `json:"id" gorm:"primaryKey,autoIncrement"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"deleted_at"`
}

func ValidatePassword(password string) error {
	passwordValidator := goval.New(goval.MinLength(6, errors.New("password cant be less than 6 characters")),
		goval.MaxLength(15, errors.New("password cant be more than 15 characters")))
	err := passwordValidator.Validate(password)
	return err
}
func validateWhiteSpaces(data interface{}) error {
	return conform.Strings(data)
}

func translateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans) + "; ")
		errs = append(errs, translatedErr)
	}
	return errs

}

type UserResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	UserResponse
	AccessToken string
}

// LoginUserToDto responsible for creating a response object for the handleLogin handler
func (u *User) LoginUserToDto(token string) *LoginResponse {
	return &LoginResponse{
		UserResponse: UserResponse{
			ID:          u.ID,
			Name:        u.Name,
			PhoneNumber: u.PhoneNumber,
			Email:       u.Email,
		},
		AccessToken: token,
	}
}

// VerifyPassword verifies the collected password with the user's hashed password
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}
