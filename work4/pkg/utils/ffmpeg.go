package utils

import (
	"work/pkg/errmsg"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func M3u8ToMp4(input, output string) error {
	err := ffmpeg.Input(input).
		Output(output, ffmpeg.KwArgs{
			"c:v":     "copy",
			"absf":    "aac_adtstoasc",
			"b:v":     "4000k",
			"bufsize": "4000k",
		}).
		OverWriteOutput().
		Run()
	if err != nil {
		return errmsg.FfmpegError
	}
	return nil
}

func GenerateMp4CoverJpg(input, output string) error {
	err := ffmpeg.Input(input).
		Output(output, ffmpeg.KwArgs{
			"ss":       "00:00:00",
			"frames:v": "1",
		}).
		OverWriteOutput().
		Run()
	if err != nil {
		return errmsg.FfmpegError
	}
	return nil
}
