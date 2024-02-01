package utils

import (
	"bytes"
	"encoding/base64"
	"image/png"

	"github.com/skip2/go-qrcode"
)

func EncodeUrlToBase64String(url string) string {
	code, _ := qrcode.New(url, qrcode.Low)
	img := code.Image(256)
	buf := bytes.NewBuffer(make([]byte, 0))
	png.Encode(buf, img)
	return `data:image/png;base64,`+base64.StdEncoding.EncodeToString(buf.Bytes())
}
