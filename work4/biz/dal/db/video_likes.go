package db

import (
	"time"
	"work/pkg/errmsg"

	"gorm.io/gorm/clause"
)

type VideoLike struct {
	Id        int64  `json:"id"`
	UserId    string `json:"user_id"`
	VideoId   string `json:"video_id"`
	CreatedAt int64  `json:"created_at"`
	DeletedAt int64  `json:"deleted_at"`
}

func GetVideoLikeList(vid string) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`video_likes`).Where(`video_id = ?`, vid).Select("user_id").Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func CreateVideoLike(videoLike *VideoLike) error {
	if err := DB.Create(videoLike).Error; err != nil {
		return errmsg.ServiceError
	}
	return nil
}

func DeleteVideoLike(vid, uid string) error {
	if err := DB.Where("video_id = ? and user_id = ?", vid, uid).Delete(&VideoLike{}).Error; err != nil {
		return errmsg.ServiceError
	}
	return nil
}

func CreateIfNotExistsVideoLike(vid string, uid string) error {
	err := DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "video_id"}, {Name: "user_id"}},
	}).Create(&VideoLike{
		UserId:    uid,
		VideoId:   vid,
		CreatedAt: time.Now().Unix(),
		DeletedAt: 0,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func GetVideoLikeListByUserId(uid string, pageNum, pageSize int64) (*[]string, error) {
	list := make([]string, 0)
	err := DB.Table(`video_likes`).Where(`user_id = ?`, uid).Select("video_id").Limit(int(pageSize)).Offset((int(pageNum-1) * int(pageSize))).Scan(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, err
}

func DeleteLikeAboutVideo(vid string) error {
	if err := DB.Where(`video_id = ?`, vid).Delete(&VideoLike{}).Error; err != nil {
		return err
	}
	return nil
}
