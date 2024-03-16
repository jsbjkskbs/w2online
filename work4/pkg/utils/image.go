package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"work/pkg/errmsg"
)

type ImageInfo struct {
}

func NewImageInfo() *ImageInfo {
	return &ImageInfo{}
}

// tag.exp: `image/jpeg`
func (ImageInfo) Get(data []byte, tag string) (height, width int, err error) {
	var imgCfg image.Config
	switch tag {
	case `image/jpeg`, `image/jpg`:
		imgCfg, err = jpeg.DecodeConfig(bytes.NewReader(data))
	case `image/png`:
		imgCfg, err = png.DecodeConfig(bytes.NewReader(data))
	default:
		return -1, -1, errmsg.FileFormatNotSupportError
	}
	if err != nil {
		return -1, -1, err
	}
	return imgCfg.Height, imgCfg.Width, err
}