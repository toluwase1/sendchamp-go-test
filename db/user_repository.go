package db

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"sendchamp-go-test/models"
)

// DB provides access to the different db

//go:generate mockgen -destination=../mocks/auth_repo_mock.go -package=mocks sendchamp-go-test/db UserRepository
type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	IsEmailExist(email string) error
	IsPhoneExist(email string) error
	FindUserByUsername(username string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	AddToBlackList(blacklist *models.BlackList) error
	TokenInBlacklist(token string) bool
	VerifyEmail(email string, token string) error
	IsTokenInBlacklist(token string) error
}

type authRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *GormDB) UserRepository {
	return &authRepo{db.DB}
}

func (a *authRepo) CreateUser(user *models.User) (*models.User, error) {
	err := a.DB.Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}
	return user, nil
}

func (a *authRepo) FindUserByUsername(username string) (*models.User, error) {
	db := a.DB
	user := &models.User{}
	err := db.Where("email = ? OR username = ?", username, username).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}
	return user, nil
}

func (a *authRepo) IsEmailExist(email string) error {
	var count int64
	err := a.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return errors.Wrap(err, "gorm.count error")
	}
	if count > 0 {
		return fmt.Errorf("email already in use")
	}
	return nil
}

func (a *authRepo) IsPhoneExist(phone string) error {
	var count int64
	err := a.DB.Model(&models.User{}).Where("phone_number = ?", phone).Count(&count).Error
	if err != nil {
		return errors.Wrap(err, "gorm.count error")
	}
	if count > 0 {
		return fmt.Errorf("phone number already in use")
	}
	return nil
}

func (a *authRepo) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := a.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *authRepo) UpdateUser(user *models.User) error {
	return nil
}

func (a *authRepo) AddToBlackList(blacklist *models.BlackList) error {
	result := a.DB.Create(blacklist)
	return result.Error
}

func (a *authRepo) TokenInBlacklist(token string) bool {
	result := a.DB.Where("token = ?", token).Find(&models.BlackList{})
	return result.Error != nil
}

func (a *authRepo) VerifyEmail(email string, token string) error {
	err := a.DB.Model(&models.User{}).Where("email = ?", email).Updates(models.User{IsEmailActive: true}).Error
	if err != nil {
		return err
	}

	err = a.AddToBlackList(&models.BlackList{Token: token})
	return err
}

func (a *authRepo) IsTokenInBlacklist(token string) error {
	var count int64
	err := a.DB.Model(&models.BlackList{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return errors.Wrap(err, "gorm.count error")
	}
	if count > 0 {
		return fmt.Errorf("token expired, request a new link")
	}
	return nil
}
