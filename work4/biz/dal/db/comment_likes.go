package db

import "time"

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

func GetCommentLikeCount(cid string) (count int64, err error) {
	if err := DB.Table(`comment_likes`).Where(`comment_id = ?`, cid).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}

func CreateCommentLike(cid, uid string) error {
	if err := DB.Create(&CommentLike{
		CommentId: cid,
		UserId:    uid,
		CreatedAt: time.Now().Unix(),
		DeletedAt: 0,
	}).Error; err != nil {
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
