package db

import (
	"fmt"
	"time"
	"work/pkg/errmsg"
)

type User struct {
	Uid       int64 `json:"uid"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatar_url"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

var UserModel User

// 创建用户信息并返回uid
func CreateUser(user *User) (string, error) {
	if err := DB.Create(user).Error; err != nil {
		return ``, err
	} else {
		return fmt.Sprint(user.Uid), err
	}
}

// 请求用户信息(username)
func QueryUserByUsername(username string) (*User, error) {
	user := User{Uid: 0}
	if err := DB.Where("username = ?", username).Find(&user).Error; err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// 请求用户信息(uid)
func QueryUserByUid(uid string) (*User, error) {
	user := User{Uid: 0}
	if err := DB.Where("uid = ?", uid).Find(&user).Error; err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// 同时用uid和username检索
func QueryUserByUidAndUsername(uid, username string) (*User, error) {
	user := User{Uid: 0}
	if err := DB.Where("uid = ? and username = ?", uid, username).Find(&user).Error; err != nil {
		return nil, errmsg.UsernameAndUidAreNotMatchedError
	} else {
		return &user, nil
	}
}

// 验证账密匹配
func VerifyUserByUsername(username, password string) (*User, error) {
	user := User{Uid: 0}
	if err := DB.Where("username = ? and password = ?", username, password).Find(&user).Error; err != nil {
		return nil, err
	}
	if user.Uid == 0 {
		return nil, errmsg.AuthenticatorError
	}
	return &user, nil
}

// 用户名判断存在
func UserIsExistByUsername(username string) (bool, error) {
	if user, err := QueryUserByUsername(username); err != nil {
		return true, err
	} else {
		return !(user.Uid == 0), nil
	}
}

// uid判断存在
func UserIsExistByUid(uid string) (bool, error) {
	if user, err := QueryUserByUid(uid); err != nil {
		return true, err
	} else {
		return !(user.Uid == 0), nil
	}
}

func UpdateAvatarUrl(uid, avatarUrl string) (*User, error) {
	exist, err := UserIsExistByUid(uid)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errmsg.ServiceError
	}

	if err = DB.Where("uid = ?", uid).Model(&UserModel).Update("avatar_url", avatarUrl).Error; err != nil {
		return nil, err
	}

	if err = DB.Where("uid = ?", uid).Model(&UserModel).Update("updated_at", time.Now().Unix()).Error; err != nil {
		return nil, err
	}

	user := User{}
	if err = DB.Where("uid = ?", uid).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}