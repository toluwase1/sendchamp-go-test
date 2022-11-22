package services

import (
	"github.com/google/uuid"
	"log"
	"net/http"
	"sendchamp-go-test/config"
	"sendchamp-go-test/db"
	"sendchamp-go-test/errors"
	apiError "sendchamp-go-test/errors"
	"sendchamp-go-test/models"
)

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks sendchamp-go-test/services TaskService
// TaskService interface
type TaskService interface {
	CreateTask(request *models.Task) (*models.Task, *apiError.Error)
	UpdateTask(request *models.Task, taskID string) *errors.Error
	DeleteTaskById(email string) *apiError.Error
	GetAllTasks(taskId string) (*models.Task, error)
}

// taskService struct
type taskService struct {
	Config   *config.Config
	taskRepo db.TaskRepository
	rabbitmq db.RabbitMqRepository
}

// NewCompanyService instantiate an taskService
func NewTaskService(taskRepo db.TaskRepository, rabbitmq db.RabbitMqRepository, conf *config.Config) TaskService {
	return &taskService{
		Config:   conf,
		taskRepo: taskRepo,
		rabbitmq: rabbitmq,
	}
}

var validTypes = map[string]bool{"medium": true, "high": true, "low": true}

func (a *taskService) CreateTask(task *models.Task) (*models.Task, *apiError.Error) {
	task.Id = uuid.New().String()
	if !validTypes[task.Priority] {
		return nil, apiError.New("invalid task type", http.StatusBadRequest)
	}
	go func() {
		err := a.rabbitmq.Rabbitmq(task)
		if err != nil {
			log.Println("Error publishing to rabbitmq", err)
		}
	}()

	task, err := a.taskRepo.CreateTask(task)
	if err != nil {
		log.Printf("task to create task: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}

	return task, nil
}

func (m *taskService) UpdateTask(request *models.Task, taskID string) *errors.Error {

	if !validTypes[request.Priority] {
		return apiError.New("invalid priority type", http.StatusBadRequest)
	}

	task := models.Task{
		Name:        request.Name,
		Description: request.Description,
		Priority:    request.Priority,
	}
	//get task where user and task id is defined above then send it for updating
	_ = m.rabbitmq.Rabbitmq(&task)
	err := m.taskRepo.UpdateTask(&task, taskID)
	if err != nil {
		return errors.ErrInternalServerError
	}
	return nil
}

func (a *taskService) DeleteTaskById(taskId string) *apiError.Error {
	err := a.taskRepo.DeleteTaskById(taskId)
	if err != nil {
		return apiError.ErrInternalServerError
	}
	task := models.Task{
		Id: taskId,
	}
	_ = a.rabbitmq.Rabbitmq(&task)
	return nil
}

func (s *taskService) GetAllTasks(taskId string) (*models.Task, error) {
	task, err := s.taskRepo.FindTaskById(taskId)
	if err != nil {
		return &models.Task{}, apiError.ErrInternalServerError
	}
	return task, nil
}
