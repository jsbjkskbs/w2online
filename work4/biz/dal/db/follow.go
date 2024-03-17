package db

import (
	"time"
)

type Follow struct {
	Id         int64  `json:"id"`
	FollowedId string `json:"followed_id"`
	FollowerId string `json:"follower_id"`
	CreatedAt  int64  `json:"created_at"`
	DeletedAt  int64  `json:"deleted_at"`
}

func CreateFollow(followedId, followerId string) error {
	err := DB.Create(&Follow{
		FollowedId: followedId,
		FollowerId: followerId,
		CreatedAt:  time.Now().Unix(),
		DeletedAt:  0,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteFollow(followedId, followerId string) error {
	if err := DB.Where("followed_id = ? and follower_id = ?", followedId, followerId).Delete(&Follow{}).Error; err != nil {
		return err
	}
	return nil
}

func GetFollowList(uid string) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`follows`).Where(`follower_id = ?`, uid).Select("followed_id").Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func GetFollowerList(uid string) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`follows`).Where(`followed_id = ?`, uid).Select("follower_id").Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func IsRelationExist(followedId, followerId string) (bool, error) {
	var count int64
	err := DB.Table(`follows`).Where(`followed_id = ? and follower_id = ?`, followedId, followerId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetFollowListPaged(uid string, pageNum, pageSize int64) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`follows`).Where(`follower_id = ?`, uid).Select(`followed_id`).Offset(int((pageNum - 1) * pageSize)).Limit(int(pageSize)).Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func GetFollowerListPaged(uid string, pageNum, pageSize int64) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.Table(`follows`).Where(`followed_id = ?`, uid).Select(`follower_id`).Offset(int((pageNum - 1) * pageSize)).Limit(int(pageSize)).Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func GetFriendListPaged(uid string, pageNum, pageSize int64) (*[]string, error) {
	list := make([]string, 0)
	if err := DB.
		Select(`follower_id`).Table(`follows`).Where(`followed_id = ? and follower_id in (
			select followed_id from follows where follower_id = ?
		)`, uid, uid).
		Offset(int((pageNum - 1) * pageSize)).Limit(int(pageSize)).Scan(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func GetFollowCount(uid string) (count int64, err error) {
	if err := DB.Table(`follows`).Where(`follower_id = ?`, uid).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}

func GetFollowerCount(uid string) (count int64, err error) {
	if err := DB.Table(`follows`).Where(`followed_id = ?`, uid).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}

func GetFriendCount(uid string) (count int64, err error) {
	if err := DB.Select(`follower_id`).Table(`follows`).Where(`followed_id = ? and follower_id in (
		select followed_id from follows where follower_id = ?
	)`, uid, uid).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}
