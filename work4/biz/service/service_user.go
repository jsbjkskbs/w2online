package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"work/biz/dal/db"
	"work/biz/model/base/user"
	"work/biz/mw/jwt"
	"work/biz/mw/redis"
	"work/pkg/constants"
	"work/pkg/errmsg"
	qiniuyunoss "work/pkg/qiniuyun_oss"
	"work/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{
		ctx: ctx,
		c:   c,
	}
}

func (service UserService) UserServiceRegister(request *user.UserRegisterRequest) (uid string, err error) {
	exist, err := db.UserIsExistByUsername(request.Username)
	if err != nil {
		return ``, err
	}
	if exist {
		return ``, errmsg.UsernameAlreadyExistError
	}
	uid, err = db.CreateUser(&db.User{
		Username:  request.Username,
		Password:  utils.EncryptBySHA256(request.Password),
		AvatarUrl: constants.DefaultAvatarUrl,
		CreatedAt: time.Now().Unix(),
		DeletedAt: 0,
		UpdatedAt: time.Now().Unix(),
	})
	if err != nil {
		return ``, err
	}
	return uid, nil
}

func (service UserService) UserServiceLogin(request *user.UserLoginRequest) (*db.User, error) {
	return db.VerifyUserByUsername(request.Username, utils.EncryptBySHA256(request.Password))
}

func (service UserService) UserServiceInfo(request *user.UserInfoRequest) (*db.User, error) {
	return db.QueryUserByUid(request.UserId)
}

func (service UserService) UserServiceAvatarUpload(request *user.UserAvatarUploadRequest) (*db.User, error) {
	uid, err := jwt.CovertJWTPayloadToString(service.ctx, service.c)
	if err != nil {
		return nil, errmsg.TokenIsInavailableError
	}

	isUploading, err := redis.IsAvatarUploading(fmt.Sprint(uid))
	if err != nil {
		return nil, errmsg.RedisError
	}
	if isUploading {
		return nil, errmsg.FileIsUploadingError
	}

	redis.AvatarSetUploadUncompleted(fmt.Sprint(uid))
	defer redis.AvatarSetUploadCompleted(fmt.Sprint(uid))

	data, err := service.uploadAvatarToOss(fmt.Sprint(uid))
	if err != nil {
		return nil, errmsg.ServiceError
	}

	return data, nil
}

func (service UserService) uploadAvatarToOss(uid string) (*db.User, error) {
	uploadRawData, err := service.c.FormFile("data")
	if err != nil {
		return nil, errmsg.FileIsUnableToBeCatchError
	}
	if uploadRawData.Size > 1*constants.MBytes {
		return nil, errmsg.FileIsTooLargeError
	}

	file, err := uploadRawData.Open()
	if err != nil {
		return nil, errmsg.ServiceError
	}
	defer file.Close()

	avatarRawData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errmsg.ServiceError
	}

	fileType := http.DetectContentType(avatarRawData)
	switch fileType {
	case `image/png`, `image/jpg`, `image/jpeg`:
		{
			var avatarUrl string
			if avatarUrl, err = qiniuyunoss.UploadAvatar(&avatarRawData, uploadRawData.Size, fmt.Sprint(uid), fileType); err != nil {
				return nil, errmsg.OssUploadError
			}
			data, err := db.UpdateAvatarUrl(uid, avatarUrl)
			if err != nil {
				return nil, errmsg.ServiceError
			}
			return data, nil
		}
	default:
		return nil, errmsg.FileFormatNotSupportError
	}
}
