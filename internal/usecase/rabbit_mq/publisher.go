package rabbit_mq

import amqp "github.com/rabbitmq/amqp091-go"

func Publish(queueName string, body []byte) error {
	_, err := ch.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	return ch.Publish(
		"", queueName, false, false,
		amqp.Publishing{ContentType: "application/octet-stream", Body: body},
	)
}
