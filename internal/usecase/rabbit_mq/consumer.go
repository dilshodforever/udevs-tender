package rabbit_mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var conn *amqp.Connection
var ch *amqp.Channel

func InitRabbitMQ() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
}

func Consume(queueName string) (<-chan amqp.Delivery, error) {
	_, err := ch.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return nil, err
	}

	return ch.Consume(queueName, "", true, false, false, false, nil)
}
