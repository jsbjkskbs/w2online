// Code generated by hertz generator.

package user

import (
	"work/biz/router/authfunc"

	"github.com/cloudwego/hertz/pkg/app"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _authMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _mfaMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _bindMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _authmfabindMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _qrcodeMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _authmfaqrcodeMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _avatarMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _uploadMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _useravataruploadMw() []app.HandlerFunc {
	// your code...
	return authfunc.Auth()
}

func _infoMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userinfoMw() []app.HandlerFunc {
	// your code...
	return authfunc.Auth()
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userloginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userregisterMw() []app.HandlerFunc {
	// your code...
	return nil
}
