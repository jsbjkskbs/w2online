package db

import (
	"fmt"
	"time"
)

type Comment struct {
	Id        int64 `json:"id"`
	UserId    string `json:"user_id"`
	VideoId   string `json:"video_id"`
	ParentId  string `json:"parent_id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

func CreateComment(comment *Comment) error {
	return DB.Create(comment).Error
}

func DeleteComment(commentId string) error {
	if err := DB.Where("id = ?", commentId).Model(&Comment{}).Update("deleted_at", time.Now().Unix()).Error; err != nil {
		return err
	}
	return nil
}

func GetChildCommentCount(commentId string) (string, error) {
	var count int64
	if err := DB.Where("parent_id = ?", commentId).Model(&Comment{}).Count(&count).Error; err != nil {
		return ``, err
	}
	return fmt.Sprint(count), nil
}

func GetVideoCommentCount(vid string) (string, error) {
	var count int64
	if err := DB.Where("video_id = ?", vid).Model(&Comment{}).Count(&count).Error; err != nil {
		return ``, err
	}
	return fmt.Sprint(count), nil
}
