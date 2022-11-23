package services

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"sendchamp-go-test/config"
	"sendchamp-go-test/db"
	"sendchamp-go-test/errors"
	apiError "sendchamp-go-test/errors"
	"sendchamp-go-test/models"
)

var (
	WsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	WsConn *websocket.Conn
)

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks sendchamp-go-test/services TaskService
// TaskService interface
type TaskService interface {
	PushTaskToDbAndSocket(task *models.Task) (*models.Task, *apiError.Error)
	PushToRabbitMq(request *models.Task) (*models.Task, *apiError.Error)
	UpdateTask(request *models.Task, taskID string) *errors.Error
	DeleteTaskById(id string) *apiError.Error
	GetTaskById(taskId string) (*models.Task, error)
	Consumer()
}

// taskService struct
type taskService struct {
	Config   *config.Config
	taskRepo db.TaskRepository
	rabbitmq db.RabbitMqRepository
}

// NewTaskService instantiate an taskService
func NewTaskService(taskRepo db.TaskRepository, rabbitmq db.RabbitMqRepository, conf *config.Config) TaskService {
	return &taskService{
		Config:   conf,
		taskRepo: taskRepo,
		rabbitmq: rabbitmq,
	}
}

var validTypes = map[string]bool{"medium": true, "high": true, "low": true}

func (a *taskService) PushToRabbitMq(task *models.Task) (*models.Task, *apiError.Error) {
	task.Id = uuid.New().String()
	if !validTypes[task.Priority] {
		return nil, apiError.New("invalid task type", http.StatusBadRequest)
	}
	go func() {
		err := a.rabbitmq.Rabbitmq(task)
		if err != nil {
			log.Println("Error publishing to rabbitmq", err)
		}
		log.Println("message successfully published in rabbitmq")
	}()
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

func (s *taskService) GetTaskById(taskId string) (*models.Task, error) {
	task, err := s.taskRepo.FindTaskById(taskId)
	if err != nil {
		return &models.Task{}, apiError.ErrInternalServerError
	}
	return task, nil
}

func (a *taskService) Consumer() {

	log.Println("consumer method")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer conn.Close()
	log.Println("Connected to RabbitMQ")
	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		panic(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume("TestQueue", "", true, false, false, false, nil)
	if err != nil {
		log.Println("Failed to register a consumer", err)
		panic(err)
	}
	log.Println("Registered a consumer")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			message, err := decodeMessage(d.Body)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("Received a message: %v \n", string(d.Body))
			log.Println("Received a message", message)
			_, errr := a.PushTaskToDbAndSocket(&message)
			if errr != nil {
				return
			}
		}
	}()

	log.Println("Waiting for messages...", msgs)
	<-forever
}

func (a *taskService) PushTaskToDbAndSocket(task *models.Task) (*models.Task, *apiError.Error) {
	task, err := a.taskRepo.CreateTask(task)
	if err != nil {
		log.Printf("task to create task: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	_, err = db.CreateServerSentEvent(task)
	if err != nil {
		return nil, apiError.New("could not send event to socket", http.StatusBadRequest)
	}
	return task, nil
}

func decodeMessage(body []byte) (models.Task, error) {
	var message models.Task
	err := json.Unmarshal(body, &message)
	if err != nil {
		log.Println(err)
		return message, err
	}
	return message, err
}
