package errmsg

import (
	"errors"
	"fmt"
)

const (
	NoErrorCode = 0

	ServiceErrCode = iota + 10000
	ParamErrCode
	AuthenticatorFailedErrCode
	RequestAlreadyExistErrCode

	UsernameAlreadyExistErrCode
	UsernameDoesNotExistErrCode
	UserDoesNotExistErrCode
	UsernameAndUidAreNotMatched

	TokenIsInavailableErrCode

	FileIsNotAImageErrCode
	ImageHeightWidthNotEqualErrCode
	FileIsUploadingErrCode

	FileIsUnableToBeCatchErrCode
	FileFormatNotSupportErrCode
	FileIsTooLargeErrCode
	FileMD5IsNotMatchErrCode

	NoSuchVideoErrCode
	NoSuchCommentErrCode

	OssUploadErrCode
	OssDeleteErrCode
	OssStatusErrCode

	RedisErrCode

	FfmpegErrCode

	ElasticErrCode

	RequestTooFreqErrCode
)

const (
	NoErrorMsg = "OK"

	ServiceErrMsg             = "Service failed"
	ParamErrMsg               = "Wrong param provided"
	AuthenticatorFailedErrMsg = "Username and password are not matched"
	RequestAlreadyExistErrMsg = "Request already exist.Perhaps send request too frequently"

	UsernameAlreadyExistErrMsg        = "This name already exists"
	UsernameDoesNotExistErrMsg        = "This name does not exist"
	UserDoesNotExistErrMsg            = "This user does not exist"
	UsernameAndUidAreNotMatchedErrMsg = "The username and uid are not matched"

	TokenIsInavailableErrMsg = "Token is inavailable"

	FileIsNotAImageErrMsg          = "Wrong image format"
	ImageHeightWidthNotEqualErrMsg = "Image's height and width are not equal"
	FileIsUploadingErrMsg          = "File is uploading"

	FileIsUnableToBeCatchErrMsg = "Can not find file uploaded"
	FileFormatNotSupportErrMsg  = "File format not support"
	FileIsTooLargeErrMsg        = "File is too large"
	FileMD5IsNotMatchErrMsg     = "Files' MD5 is not matched"

	NoSuchVideoErrMsg   = "The video doesn't exist"
	NoSuchCommentErrMsg = "The comment doesn't exist"

	OssUploadErrMsg = "Oss upload failed"
	OssDeleteErrMsg = "Oss delete failed"
	OssStatusErrMsg = "Oss get file info failed"

	RedisErrMsg = "Redis Error"

	FfmpegErrMsg = "Ffmpeg Error"

	ElasticErrMsg = "Elastic Error"

	RequestTooFreqErrMsg = "too many request, pls wait"
)

type ErrorMessage struct {
	ErrorCode int64
	ErrorMsg  string
}

func (err ErrorMessage) Error() string {
	return fmt.Sprintf("%v, Code:%v", err.ErrorMsg, err.ErrorCode)
}

func (err ErrorMessage) WithMessage(msg string) ErrorMessage {
	return ErrorMessage{
		ErrorCode: err.ErrorCode,
		ErrorMsg:  msg,
	}
}

func NewErrorMessage(code int64, msg string) ErrorMessage {
	return ErrorMessage{
		ErrorCode: code,
		ErrorMsg:  msg,
	}
}

var (
	NoError = NewErrorMessage(NoErrorCode, NoErrorMsg)

	ServiceError             = NewErrorMessage(ServiceErrCode, ServiceErrMsg)
	ParamError               = NewErrorMessage(ParamErrCode, ParamErrMsg)
	AuthenticatorError       = NewErrorMessage(AuthenticatorFailedErrCode, AuthenticatorFailedErrMsg)
	RequestAlreadyExistError = NewErrorMessage(RequestAlreadyExistErrCode, RequestAlreadyExistErrMsg)

	UsernameAlreadyExistError        = NewErrorMessage(UsernameAlreadyExistErrCode, UsernameAlreadyExistErrMsg)
	UsernameDoesNotExistError        = NewErrorMessage(UsernameDoesNotExistErrCode, UsernameDoesNotExistErrMsg)
	UserDoesNotExistError            = NewErrorMessage(UserDoesNotExistErrCode, UserDoesNotExistErrMsg)
	UsernameAndUidAreNotMatchedError = NewErrorMessage(UsernameAndUidAreNotMatched, UsernameAndUidAreNotMatchedErrMsg)

	TokenIsInavailableError = NewErrorMessage(TokenIsInavailableErrCode, TokenIsInavailableErrMsg)

	FileIsNotAImageError          = NewErrorMessage(FileIsNotAImageErrCode, FileIsNotAImageErrMsg)
	ImageHeightWidthNotEqualError = NewErrorMessage(ImageHeightWidthNotEqualErrCode, ImageHeightWidthNotEqualErrMsg)
	FileIsUploadingError          = NewErrorMessage(FileIsUploadingErrCode, FileIsUploadingErrMsg)

	FileIsUnableToBeCatchError = NewErrorMessage(FileIsUnableToBeCatchErrCode, FileIsUnableToBeCatchErrMsg)
	FileFormatNotSupportError  = NewErrorMessage(FileFormatNotSupportErrCode, FileFormatNotSupportErrMsg)
	FileIsTooLargeError        = NewErrorMessage(FileIsTooLargeErrCode, FileIsTooLargeErrMsg)
	FileMD5IsNotMatchError     = NewErrorMessage(FileMD5IsNotMatchErrCode, FileMD5IsNotMatchErrMsg)

	NoSuchVideoError   = NewErrorMessage(NoSuchVideoErrCode, NoSuchVideoErrMsg)
	NoSuchCommentError = NewErrorMessage(NoSuchCommentErrCode, NoSuchCommentErrMsg)

	OssUploadError = NewErrorMessage(OssUploadErrCode, OssUploadErrMsg)
	OssDeleteError = NewErrorMessage(OssDeleteErrCode, OssDeleteErrMsg)
	OssStatusError = NewErrorMessage(OssStatusErrCode, OssStatusErrMsg)

	RedisError = NewErrorMessage(RedisErrCode, RedisErrMsg)

	FfmpegError = NewErrorMessage(FfmpegErrCode, FfmpegErrMsg)

	ElasticError = NewErrorMessage(ElasticErrCode, ElasticErrMsg)

	RequestTooFreqError = NewErrorMessage(RequestTooFreqErrCode, RequestTooFreqErrMsg)
)

// 提供转换方法
func Convert(err error) ErrorMessage {
	var e ErrorMessage
	if errors.As(err, &e) {
		return e
	}

	s := ServiceError
	s.ErrorMsg = err.Error()
	return s
}
