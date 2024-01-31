package redis

import (
	"encoding/json"
	"work/pkg/utils"
)

type _Message struct {
	From      string `json:"f"`
	Content   string `json:"c"`
	Timestamp int64  `json:"t"`
}

func PushMessage(from, to, content string, timestamp int64) error {
	data, err := json.Marshal(&_Message{
		From:      from,
		Content:   content,
		Timestamp: timestamp,
	})
	if err != nil {
		return err
	}
	if _, err := redisDBChatInfo.RPush(to, string(data)).Result(); err != nil {
		return err
	}
	return nil
}

func IsMessageQueueEmpty(to string) bool {
	exist, err := redisDBChatInfo.Exists(to).Result()
	if err != nil {
		return true
	}
	return exist == 0
}

// Returns are `from`,`content`,`when`,`err`
func PopMessage(to string) (string, string, string, error) {
	data, err := redisDBChatInfo.LPop(to).Result()
	if err != nil {
		return ``, ``, ``, err
	}
	var msg _Message
	if err = json.Unmarshal([]byte(data), &msg); err != nil {
		return ``, ``, ``, err
	}
	return msg.From, msg.Content, utils.ConvertTimestampToStringDefault(msg.Timestamp), nil
}
