package db

import (
	"time"

	"gorm.io/gorm/clause"
)

type CommentLike struct {
	Id        int64  `json:"id"`
	UserId    string `json:"user_id"`
	CommentId string `json:"comment_id"`
	CreatedAt int64  `json:"created_at"`
	DeletedAt int64  `json:"deleted_at"`
}

func GetCommentLikeList(cid string) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`comment_likes`).Where(`comment_id = ?`, cid).Select(`user_id`).Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func CreateCommentLike(commentLike *CommentLike) error {
	if err := DB.Create(commentLike).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCommentLike(cid, uid string) error {
	if err := DB.Where(`comment_id = ? and user_id = ?`, cid, uid).Delete(&CommentLike{}).Error; err != nil {
		return err
	}
	return nil
}

func CreateIfNotExistsCommentLike(cid, uid string) error {
	err := DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "video_id"}, {Name: "user_id"}},
	}).Create(&CommentLike{
		UserId:    uid,
		CommentId: cid,
		CreatedAt: time.Now().Unix(),
		DeletedAt: 0,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteCommentLikeAboutComment(cid string) error {
	if err := DB.Where(`comment_id = ?`, cid).Delete(&CommentLike{}).Error; err != nil {
		return err
	}
	return nil
}
