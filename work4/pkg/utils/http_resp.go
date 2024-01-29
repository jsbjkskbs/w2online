package utils

import (
	"errors"
	"work/pkg/errmsg"
)

type BaseHttpResponse struct {
	StatusCode int64
	StatusMsg  string
}

func baseHttpResponse(err errmsg.ErrorMessage) *BaseHttpResponse {
	return &BaseHttpResponse{
		StatusCode: err.ErrorCode,
		StatusMsg:  err.ErrorMsg,
	}
}

func CreateBaseHttpResponse(err error) *BaseHttpResponse {
	if err == nil {
		return baseHttpResponse(errmsg.NoError)
	}

	e:=errmsg.ErrorMessage{}
	if errors.As(err,&e){
		return baseHttpResponse(e)
	}

	s:=errmsg.ServiceError.WithMessage(err.Error())
	return baseHttpResponse(s)
}
