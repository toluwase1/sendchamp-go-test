package db

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"sendchamp-go-test/models"
)

// DB provides access to the different db

//go:generate mockgen -destination=../mocks/auth_repo_mock.go -package=mocks sendchamp-go-test/db AuthRepository
type TaskRepository interface {
	CreateTask(task *models.Task) (*models.Task, error)
	FindTaskByName(email string) (*models.Task, error)
	UpdateTask(task *models.Task, taskId string) error
	DeleteTaskById(taskId string) error
	FindTaskById(email string) (*models.Task, error)
	IsTaskOwner(taskId string, userId float64) error
}

type taskRepo struct {
	DB *gorm.DB
}

func NewTaskRepo(db *GormDB) TaskRepository {
	return &taskRepo{db.DB}
}

func (a *taskRepo) CreateTask(task *models.Task) (*models.Task, error) {
	err := a.DB.Create(task).Error
	if err != nil {
		return nil, fmt.Errorf("could not create task: %v", err)
	}
	return task, nil
}

func (a *taskRepo) FindTaskByName(username string) (*models.Task, error) {
	db := a.DB
	user := &models.Task{}
	err := db.Where("email = ? OR username = ?", username, username).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}
	return user, nil
}

func (m *taskRepo) UpdateTask(task *models.Task, taskId string) error {
	err := m.DB.Model(&models.Task{}).
		Where("id = ?", taskId).
		Updates(task).Error
	if err != nil {
		return fmt.Errorf("could not update task: %v", err)
	}
	return nil
}

func (a *taskRepo) DeleteTaskById(taskId string) error {
	task := &models.Task{}
	err := a.DB.Where("id = ?", taskId).Find(task).Error
	if err != nil {
		return fmt.Errorf("could not find task to delete: %v", err)
	}
	err = a.DB.Delete(&models.Task{}, "id = ?", taskId).Error
	if err != nil {
		return fmt.Errorf("could not delete task: %v", err)
	}
	return nil
}

func (a *taskRepo) FindTaskById(id string) (*models.Task, error) {
	task := &models.Task{}
	err := a.DB.Where("id = ?", id).Find(task).Error
	if err != nil {
		return &models.Task{}, fmt.Errorf("could not find any task with that id %v", err)
	}
	return task, nil
}

func (a *taskRepo) IsTaskOwner(taskId string, userId float64) error {
	var count int64
	err := a.DB.Model(&models.Task{}).Where("id = ? AND user_id = ?", taskId, userId).Count(&count).Error
	if err != nil {
		return errors.Wrap(err, "gorm.count error")
	}
	if count < 1 {
		return fmt.Errorf("not authorized")
	}
	return nil
}
