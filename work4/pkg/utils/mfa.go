package utils

import (
	"fmt"

	"github.com/pquerna/otp/totp"
)

type AuthController struct {
	uid    string
	code   string
	secret string
}

type MfaAuthInfo struct {
	Url    string
	Secret string
}

func NewAuthController(uid, code, secret string) *AuthController {
	return &AuthController{uid: uid, code: code, secret: secret}
}

func (ac AuthController) GenerateTOTP() (*MfaAuthInfo, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      `127.0.0.1`,
		AccountName: fmt.Sprint(ac.uid),
	})

	if err != nil {
		return nil, err
	}

	return &MfaAuthInfo{Url: key.URL(), Secret: key.Secret()}, nil
}

func (ac AuthController) VerifyTOTP() bool {
	valid := totp.Validate(ac.code, ac.secret)
	return valid
}
