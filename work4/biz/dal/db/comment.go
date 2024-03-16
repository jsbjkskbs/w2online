package db

import (
	"fmt"
)

type Comment struct {
	Id        int64  `json:"id"`
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
	if err := DB.Where("id = ?", commentId).Delete(&Comment{}).Error; err != nil {
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

func GetVideoCommentList(vid string) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`comments`).Where(`video_id = ?`, vid).Select("id").Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func GetCommentChildList(cid string) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`comments`).Where(`parent_id = ?`, cid).Select(`id`).Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func GetCommentIdList() (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`comments`).Select("id").Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func DeleteCommentAndCommentLikeAboutVideo(vid string) error {
	list, err := GetVideoCommentList(vid)
	if err != nil {
		return err
	}
	for _, item := range *list {
		if err := DeleteCommentLikeAboutComment(item); err != nil {
			return err
		}
	}
	if err := DB.Where(`video_id = ?`, vid).Delete(&Comment{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteChildComment(cid string) error {
	if err := DB.Where("parent_id = ?", cid).Delete(&Comment{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteChildAndLikesOfParentAndChild(cid string) error {
	list, err := GetCommentChildList(cid)
	if err != nil {
		return err
	}
	if err := DeleteChildComment(cid); err != nil {
		return err
	}
	if err := DeleteCommentLikeAboutComment(cid); err != nil {
		return err
	}
	for _, item := range *list {
		if err := DeleteCommentLikeAboutComment(item); err != nil {
			return err
		}
	}
	return nil
}
