package dustman

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"
	"work/biz/mw/redis"
	"work/pkg/errmsg"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type RedisVideoDustman struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRedisDustman() *RedisVideoDustman {
	ctx, cancle := context.WithCancel(context.Background())
	return &RedisVideoDustman{
		ctx:    ctx,
		cancel: cancle,
	}
}

func (rs RedisVideoDustman) Run() {
	go func() {
		for {
			select {
			case <-rs.ctx.Done():
				hlog.Info("OK,RedisDustman stop running.")
				return
			default:
				time.Sleep(1 * time.Hour)
				keys, err := redis.GetVideoDBKeys()
				if err != nil {
					hlog.Warn(errmsg.RedisErrMsg)
				}
				timestampNow := time.Now().Unix()
				keysDel := make([]string, 0)
				for _, key := range keys {
					if timestampNow-getKeyTimetamp(key) > 86400000 { //超过一天
						keysDel = append(keysDel, key)
					}
				}
				redis.DelVideoDBKeys(keysDel)
			}
		}
	}()
}

func (rs RedisVideoDustman) Stop() {
	rs.cancel()
}

func getKeyTimetamp(key string) int64 {
	info := strings.Split(key, `:`)
	timestamp, _ := strconv.ParseInt(info[2], 10, 64)
	return timestamp
}

type FileVideoDustman struct {
	ctx    context.Context
	cancle context.CancelFunc
}

func NewFileDustman() *FileVideoDustman {
	ctx, cancle := context.WithCancel(context.Background())
	return &FileVideoDustman{
		ctx:    ctx,
		cancle: cancle,
	}
}

const (
	fileVideoBlacklist = "./pkg/data/temp/video"
)

func (fd FileVideoDustman) Run() {
	go func() {
		for {
			select {
			case <-fd.ctx.Done():
				hlog.Info("OK,FileDustman stop running.")
				return
			default:
				time.Sleep(1 * time.Hour)
				foldernames, err := getFoldernames(fileVideoBlacklist)
				if err != nil {
					hlog.Warn("Dustman can not get foldernames")
				}
				timestampNow := time.Now().Unix()
				for _, foldername := range foldernames {
					if timestampNow-getFolderTimestamp(foldername) > 86400000 {
						os.RemoveAll(fileVideoBlacklist + `/` + foldername)
					}
				}
			}
		}
	}()
}

func (fd FileVideoDustman) Stop() {
	fd.cancle()
}

func getFoldernames(rootpath string) ([]string, error) {
	folders, err := os.ReadDir(rootpath)
	if err != nil {
		return nil, err
	}
	foldernames := make([]string, 0)
	for _, folder := range folders {
		foldernames = append(foldernames, folder.Name())
	}
	return foldernames, nil
}

func getFolderTimestamp(foldername string) int64 {
	info := strings.Split(foldername, `_`)
	timestamp, _ := strconv.ParseInt(info[1], 10, 64)
	return timestamp
}
