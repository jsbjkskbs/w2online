package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"work/pkg/errmsg"
)

func GetFileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return ``, errmsg.ServiceError
	}
	defer file.Close()

	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, file); err != nil {
		return ``, errmsg.ServiceError
	}

	return fmt.Sprint(hex.EncodeToString(md5Hash.Sum(nil))), nil
}

func GetBytesMD5(data []byte) string {
	hash:=md5.Sum(data)
	return hex.EncodeToString(hash[:])
}