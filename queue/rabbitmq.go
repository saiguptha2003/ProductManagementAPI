package queue

import (
	"log"

	"github.com/streadway/amqp"
)

var Channel *amqp.Channel

func InitQueue() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}

	Channel, err = conn.Channel()
	if err != nil {
		return err
	}

	_, err = Channel.QueueDeclare(
		"image_processing_queue",
		true, 
		false, 
		false,
		false, 
		nil,   
	)
	if err != nil {
		return err
	}

	log.Println("RabbitMQ initialized successfully")
	return nil
}

func Publish(message []byte) error {
	return Channel.Publish(
		"",
		"image_processing_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
