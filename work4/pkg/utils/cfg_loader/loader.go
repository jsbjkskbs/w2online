package cfgloader

import (
	"work/biz/dal"
	"work/biz/mw/elasticsearch"
	"work/biz/mw/redis"
	"work/pkg/constants"
	qiniuyunoss "work/pkg/qiniuyun_oss"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var globalConfig *viper.Viper

func Run() error {
	globalConfig = viper.New()
	globalConfig.SetConfigName("config")
	globalConfig.SetConfigType("json")
	globalConfig.AddConfigPath(configPath)
	if err := globalConfig.ReadInConfig(); err != nil {
		return err
	}
	loadConfig()

	go func() {
		globalConfig.WatchConfig()
		globalConfig.OnConfigChange(func(e fsnotify.Event) {
			hlog.Info("modifying the cfg file has detected")
			loadConfig()
		})
	}()

	return nil
}

func loadConfig() {
	constants.MysqlDSN = globalConfig.GetString("MysqlDSN")
	dal.Load()

	constants.RedisAddr = globalConfig.GetString("RedisAddr")
	constants.RedisPassword = globalConfig.GetString("RedisPassword")
	redis.Load()

	constants.ElasticAddr = globalConfig.GetString("ElasticAddr")
	elasticsearch.Load()

	qiniuyunoss.Bucket = globalConfig.GetString("OssBucket")
	qiniuyunoss.SecretKey = globalConfig.GetString("OssSecretKey")
	qiniuyunoss.AccessKey = globalConfig.GetString("OssAccessKey")
	qiniuyunoss.Url = globalConfig.GetString("OssUrl")
	qiniuyunoss.Load()
}
