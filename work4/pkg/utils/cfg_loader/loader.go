package cfgloader

import (
	"work/biz/dal"
	"work/biz/mw/elasticsearch"
	"work/biz/mw/rabbitmq"
	"work/biz/mw/redis"
	"work/biz/mw/sentinel"
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
	hlog.Info("MysqlDSN:" + constants.MysqlDSN)
	dal.Load()

	constants.RedisAddr = globalConfig.GetString("RedisAddr")
	constants.RedisPassword = globalConfig.GetString("RedisPassword")
	hlog.Info("RedisAddr:" + constants.RedisAddr)
	redis.Load()

	constants.ElasticAddr = globalConfig.GetString("ElasticAddr")
	hlog.Info("ElasticAddr:" + constants.ElasticAddr)
	elasticsearch.Load()

	constants.RabbitmqDSN = globalConfig.GetString("RabbitmqDSN")
	hlog.Info("RabbitmqDSN:" + constants.RabbitmqDSN)
	rabbitmq.Load()

	qiniuyunoss.Bucket = globalConfig.GetString("OssBucket")
	qiniuyunoss.SecretKey = globalConfig.GetString("OssSecretKey")
	qiniuyunoss.AccessKey = globalConfig.GetString("OssAccessKey")
	qiniuyunoss.Url = globalConfig.GetString("OssUrl")
	qiniuyunoss.Load()

	sentinel.Rules = globalConfig.GetStringMap("SentinelRules")
	sentinel.Load()
}
