package db

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"log"
	"sendchamp-go-test/models"
)

type RabbitMqRepository interface {
	Rabbitmq(msg *models.Task) error
}

type rabbitRepo struct {
	DB *gorm.DB
}

func NewRabbitRepo(db *GormDB) RabbitMqRepository {
	return &rabbitRepo{db.DB}
}
func (k *rabbitRepo) Rabbitmq(msg *models.Task) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	fmt.Println("Connected to RabbitMQ")
	ch, err := conn.Channel()

	if err != nil {
		log.Println("Failed to open a channel")
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("TestQueue", false, false, false, false, nil)
	if err != nil {
		log.Println("Failed to declare a queue")
		panic(err)
	}
	fmt.Println(q)
	msgBody, err := json.Marshal(msg)
	if err != nil {
		log.Println("Failed to marshal message")
		panic(err)
	}

	err = ch.Publish("", "TestQueue", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(msgBody),
	})
	return err
}
