package rabbitmq

import (
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func (r RabbitMq) Subscribe(pChan *chan []byte) {
	_, err := r.Chan.QueueDeclare(
		r.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		hlog.Error(err)
	}

	messageChan, err := r.Chan.Consume(
		r.QueueName,
		fmt.Sprint(time.Now().Unix()),
		true,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		hlog.Error(err)
	}

	go func() {
		for message := range messageChan {
			*pChan <- message.Body
		}
	}()
}
