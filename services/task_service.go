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

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService
// CompanyService interface
type CompanyService interface {
	CreateTask(request *models.Task) (*models.Task, *apiError.Error)
	UpdateTask(request *models.Task, companyID string) *errors.Error
	DeleteTaskById(email string) *apiError.Error
	GetAllTasks(companyId string) (*models.Task, error)
}

// taskService struct
type taskService struct {
	Config   *config.Config
	taskRepo db.TaskRepository
	kafka    db.RabbitMqRepository
}

// NewCompanyService instantiate an taskService
func NewTaskService(taskRepo db.TaskRepository, kafkaRepo db.RabbitRepository, conf *config.Config) CompanyService {
	return &taskService{
		Config:   conf,
		taskRepo: taskRepo,
		kafka:    kafkaRepo,
	}
}

var validTypes = map[string]bool{"medium": true, "high": true, "low": true}

func (a *taskService) CreateTask(task *models.Task) (*models.Task, *apiError.Error) {
	task.Id = uuid.New().String()
	if !validTypes[task.Priority] {
		return nil, apiError.New("invalid task type", http.StatusBadRequest)
	}
	task, err := a.taskRepo.CreateTask(task)
	if err != nil {
		log.Printf("task to create task: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	_, _ = a.kafka.AddMessageToRabbitMq(task, context.Background())

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
	err := m.companyRepo.UpdateCompany(&task, taskID)
	if err != nil {
		return errors.ErrInternalServerError
	}
	_, _ = m.kafka.AddMessageToRabbitMq(&task, context.Background())
	return nil
}

func (a *taskService) DeleteTaskById(companyId string) *apiError.Error {
	err := a.companyRepo.DeleteCompanyById(companyId)
	if err != nil {
		return apiError.ErrInternalServerError
	}
	coy := models.Task{
		Id: companyId,
	}
	_, _ = a.kafka.AddMessageToRabbitMq(&coy, context.Background())
	return nil
}

func (s *taskService) GetAllTasks(companyId string) (*models.Task, error) {
	coy, err := s.companyRepo.FindCompanyById(companyId)
	if err != nil {
		return &models.Task{}, apiError.ErrInternalServerError
	}
	return coy, nil
}
