package rabbitmq

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"work/biz/dal/db"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type CommentQueueStruct struct {
	publisher  *RabbitMq
	subscriber *RabbitMq
	dataChan   chan db.Comment
	msgChan    chan []byte
	ctx        context.Context
	cancle     context.CancelFunc
	sysCtx     context.Context
	sysCancle  context.CancelFunc
}

func NewCommentMQ() *CommentQueueStruct {
	const (
		queueName    = `commentQueue`
		exchangeName = `commentExchange`
		routingKey   = `commentKey`
	)
	publisher, err := NewRabbitMq(queueName, exchangeName, routingKey)
	if err != nil {
		return nil
	}

	subscriber, err := NewRabbitMq(queueName, exchangeName, routingKey)
	if err != nil {
		return nil
	}
	ctx, cancle := context.WithCancel(context.Background())
	sysctx, sysCancle := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	return &CommentQueueStruct{
		publisher:  publisher,
		subscriber: subscriber,
		dataChan:   make(chan db.Comment, 65536),
		msgChan:    make(chan []byte, 65536),
		ctx:        ctx,
		cancle:     cancle,
		sysCtx:     sysctx,
		sysCancle:  sysCancle,
	}
}

func (mq CommentQueueStruct) Run() {
	mq.subscriber.Subscribe(&mq.msgChan)

	go func() {
		for {
			select {
			case <-mq.sysCtx.Done():
				hlog.Warn("suddenly stop CommentMQ may cause data loss.")
				return
			case <-mq.ctx.Done():
				hlog.Warn("suddenly stop CommentMQ may cause data loss.")
				return
			case msg := <-mq.msgChan:
				comment := db.Comment{}
				if err := json.Unmarshal(msg, &comment); err != nil {
					hlog.Error(err)
					continue
				}
				if err := db.CreateComment(&comment); err != nil {
					hlog.Error(err)
				}
			}
		}
	}()
}

func (mq CommentQueueStruct) Send(data *db.Comment) error {
	return mq.publisher.Publish(*data)
}

func (mq CommentQueueStruct) Stop() {
	mq.cancle()
	mq.publisher.ReleaseRabbitMq()
	mq.subscriber.ReleaseRabbitMq()
}
