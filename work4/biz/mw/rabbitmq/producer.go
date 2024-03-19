package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

func (r RabbitMq) Publish(data interface{}) error {
	_, err := r.Chan.QueueDeclare(
		r.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = r.Chan.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = r.Chan.QueueBind(
		r.QueueName,
		r.RoutingKey,
		r.Exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	message, _ := json.Marshal(data)
	return r.Chan.Publish(
		r.Exchange,
		r.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
