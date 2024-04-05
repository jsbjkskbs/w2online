package db

import "fmt"

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

func GetChildCommentCount(commentId string) (int64, error) {
	var count int64
	if err := DB.Where("parent_id = ?", commentId).Model(&Comment{}).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}

func GetVideoCommentCount(vid string) (string, error) {
	var count int64
	if err := DB.Where("video_id = ?", vid).Model(&Comment{}).Count(&count).Error; err != nil {
		return ``, err
	}
	return fmt.Sprint(count), nil
}

func GetCommentInfo(commentId string) (*Comment, error) {
	var comment Comment
	if err := DB.Table(`comments`).Where(`id = ?`, commentId).Find(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func GetParentCommentId(commentId string) (string, error) {
	var parentId string
	if err := DB.Table(`comments`).Where(`id = ?`, commentId).Select(`parent_id`).Find(&parentId).Error; err != nil {
		return ``, err
	}
	return parentId, nil
}

func GetCommentVideoId(commentId string) (string, error) {
	var videoId string
	if err := DB.Table(`comments`).Where(`id = ?`, commentId).Select(`video_id`).Find(&videoId).Error; err != nil {
		return ``, err
	}
	return videoId, nil
}

func GetVideoCommentList(vid string) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`comments`).Where(`video_id = ?`, vid).Select("id").Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func GetVideoCommentListByPart(vid string, pageNum, pageSize int64) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`comments`).Where(`video_id = ?`, vid).Select("id").Limit(int(pageSize)).Offset(int(pageNum-1) * int(pageSize)).Scan(&list).Error; err != nil {
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

func GetCommentChildListByPart(cid string, pageNum, pageSize int64) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`comments`).Where(`parent_id = ?`, cid).Select(`id`).Limit(int(pageSize)).Offset(int(pageNum-1) * int(pageSize)).Scan(&list).Error; err != nil {
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

func IsCommentExist(cid string) (bool, error) {
	var count int64
	if err := DB.Table(`comments`).Where(`id = ?`, cid).Count(&count).Error; err != nil {
		return false, err
	}
	return count != 0, nil
}
