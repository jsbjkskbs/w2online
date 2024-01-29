package db

import (
	"fmt"
	"work/pkg/errmsg"
)

type Video struct {
	Id          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	VideoUrl    string `json:"video_url"`
	CoverUrl    string `json:"cover_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	VisitCount  int64  `json:"visit_count"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	DeletedAt   int64  `json:"deleted_at"`
}

func CreateVideo(video *Video) (string, error) {
	if err := DB.Create(video).Error; err != nil {
		return ``, err
	} else {
		return fmt.Sprint(video.Id), nil
	}
}

func UpdateVideoUrl(videoUrl, coverUrl, vid string) error {
	if err := DB.Where("id = ?", vid).Model(&Video{}).Update("video_url", videoUrl).Error; err != nil {
		return errmsg.ServiceError
	}
	if err := DB.Where("id = ?", vid).Model(&Video{}).Update("cover_url", coverUrl).Error; err != nil {
		return errmsg.ServiceError
	}
	return nil
}

func UpdateVideoVisit(vid, visitCount string) error {
	if err := DB.Where("id = ?", vid).Model(&Video{}).Update("visit_count", visitCount); err != nil {
		return errmsg.ServiceError
	}
	return nil
}

func GetVideoIdList() (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`videos`).Select("id").Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func GetVideoVisitCount(vid string) (count string, err error) {
	if err = DB.Table(`videos`).Select(`visit_count`).Where(`id = ?`, vid).Scan(&count).Error; err != nil {
		return ``, err
	}
	return count, err
}
