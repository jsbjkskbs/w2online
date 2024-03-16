package qiniuyunoss

import (
	"bytes"
	"context"
	"os"
	"work/pkg/errmsg"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

/*

	注:oss的密钥一类的涉密参数保存在同文件夹内的cfg.go，不会同步到github上
	bucket|url|secretKey|accessKey
*/

var (
	formUploader   *storage.FormUploader
	bucketManager  *storage.BucketManager
	resumeUploader *storage.ResumeUploaderV2
	mac            *auth.Credentials
	upToken        string
)

func Load() {
	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac = qbox.NewMac(AccessKey, SecretKey)
	upToken = putPolicy.UploadToken(mac)
	cfg := storage.Config{
		UseHTTPS:      false,
		UseCdnDomains: false,

		Zone: &storage.ZoneHuanan, //your bucket zone
	}
	formUploader = storage.NewFormUploader(&cfg)
	bucketManager = storage.NewBucketManager(mac, &cfg)
	resumeUploader = storage.NewResumeUploaderV2(&cfg)

	hlog.Info("Oss service prepared successfully")
}

func UploadAvatar(data *[]byte, dataSize int64, uid string, tag string) (string, error) {
	deleteAvatar(uid)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	var suffix string
	switch tag {
	case `image/jpeg`, `image/jpg`:
		suffix = `.jpg`
	case `image/png`:
		suffix = `.png`
	}
	err := formUploader.Put(
		context.Background(),
		&ret,
		upToken,
		`avatar/`+uid+suffix,
		bytes.NewReader(*data),
		dataSize,
		&putExtra,
	)
	if err != nil {
		return ``, errmsg.OssUploadError
	}
	return Url + `/avatar/` + uid + suffix, nil
}

func deleteAvatar(uid string) {
	keys := []string{
		`avatar/` + uid + `.jpg`,
		`avatar/` + uid + `.jpeg`,
		`avatar/` + uid + `.png`,
	}
	for _, key := range keys {
		bucketManager.Delete(Bucket, key)
	}
}

func UploadVideo(path, vid string) (string, error) {
	recoder, err := storage.NewFileRecorder(os.TempDir())
	if err != nil {
		return ``, errmsg.OssUploadError
	}
	ret := storage.PutRet{}
	putExtra := storage.RputV2Extra{
		Recorder: recoder,
	}
	err = resumeUploader.PutFile(
		context.Background(),
		&ret,
		upToken,
		`video/`+vid+`/video.mp4`,
		path,
		&putExtra,
	)
	if err != nil {
		return ``, errmsg.OssUploadError
	}
	return Url + `/video/` + vid + `/video.mp4`, nil
}

func UploadVideoCover(path, vid string) (string, error) {
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	err := formUploader.PutFile(
		context.Background(),
		&ret,
		upToken,
		`video/`+vid+`/cover.jpg`,
		path,
		&putExtra,
	)
	if err != nil {
		return ``, errmsg.OssUploadError
	}
	return Url + `/video/` + vid + `/cover.jpg`, nil
}
