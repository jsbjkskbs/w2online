package db

import "time"

type Message struct {
	Id         int64  `json:"id"`
	FromUserId string `json:"from_user_id"`
	ToUserId   string `json:"to_user_id"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
	DeletedAt  int64  `json:"deleted_at"`
}

func CreateMessage(from, to, content string) error {
	if err := DB.Create(&Message{
		FromUserId: from,
		ToUserId:   to,
		Content:    content,
		CreatedAt:  time.Now().Unix(),
		DeletedAt:  0,
	}).Error; err != nil {
		return err
	}
	return nil
}

func PopMessage(to string) (*[]Message, error) {
	list := make([]Message, 0)
	err := DB.Table(`messages`).Where(`to_user_id = ?`, to).Scan(&list).Error
	if err != nil {
		return nil, err
	}
	err = DB.Table(`messages`).Where(`to_user_id = ?`, to).Delete(&Message{}).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}
