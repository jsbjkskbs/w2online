package rabbitmq

import (
	"work/pkg/constants"

	"github.com/streadway/amqp"
)

type RabbitMq struct {
	Conn       *amqp.Connection
	Chan       *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string
}

func NewRabbitMq(queueName, exchange, routingkey string) (*RabbitMq, error) {
	rabbitMq := RabbitMq{
		QueueName:  queueName,
		Exchange:   exchange,
		RoutingKey: routingkey,
	}

	var err error
	if rabbitMq.Conn, err = amqp.Dial(constants.RabbitmqDSN); err != nil {
		return nil, err
	}

	if rabbitMq.Chan, err = rabbitMq.Conn.Channel(); err != nil {
		return nil, err
	}

	return &rabbitMq, nil
}

func (r RabbitMq) ReleaseRabbitMq() {
	r.Conn.Close()
	r.Chan.Close()
}
