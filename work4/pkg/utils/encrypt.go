package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func EncryptBySHA256(password string) string {
	return hex.EncodeToString(sha256.New().Sum([]byte(password)))
}
