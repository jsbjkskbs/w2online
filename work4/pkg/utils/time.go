package utils

import "time"

const defaultTimeFormat = `2006-01-02 15:04:05`

var timeLocation, _ = time.LoadLocation("Asia/Shanghai")

//<<中国的时区为什么是Asia/Shanghai，而不是Asia/Beijing？>> https://zhuanlan.zhihu.com/p/450867597

func ConvertTimestampToStringDefault(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(defaultTimeFormat)
}

func CovertStringToTimestampDefault(date string) (int64, error) {
	t, err := time.ParseInLocation(defaultTimeFormat, date, timeLocation)
	return t.Unix(), err
}

func ConvertTimestampToString(timestamp int64, fomat string) string {
	return time.Unix(timestamp, 0).Format(fomat)
}
