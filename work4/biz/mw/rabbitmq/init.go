package rabbitmq

var (
	CommentMQ *CommentQueueStruct = nil
)

func Load() {
	if CommentMQ != nil {
		CommentMQ.Stop()
	}

	CommentMQ = NewCommentMQ()
	CommentMQ.Run()
}
