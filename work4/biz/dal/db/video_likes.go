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
	if err := DB.Where("video_id = ? and user_id = ?", vid, uid).Model(&VideoLike{}).Update("deleted_at", time.Now().Unix()).Error; err != nil {
		return errmsg.ServiceError
	}
	return nil
}

func CreateIfNotExistsVideoLike(vid string, likeList *[]string) error {
	for _, uid := range *likeList {
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
	}
	return nil
}
