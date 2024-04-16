package cfgloader

import (
	"work/biz/dal"
	"work/biz/mw/elasticsearch"
	"work/biz/mw/rabbitmq"
	"work/biz/mw/redis"
	"work/biz/mw/sentinel"
	"work/biz/service"
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
	globalConfig.SetConfigType("yaml")
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

	redisConfig:=globalConfig.GetStringMapString("Redis")
	constants.RedisAddr = redisConfig["address"]
	constants.RedisPassword = redisConfig["password"]
	hlog.Info("RedisAddr:" + constants.RedisAddr)
	redis.Load()

	constants.ElasticAddr = globalConfig.GetString("ElasticAddr")
	hlog.Info("ElasticAddr:" + constants.ElasticAddr)
	elasticsearch.Load()

	constants.RabbitmqDSN = globalConfig.GetString("RabbitmqDSN")
	hlog.Info("RabbitmqDSN:" + constants.RabbitmqDSN)
	rabbitmq.Load()

	qiniuyunConfig := globalConfig.GetStringMapString("OSS")
	qiniuyunoss.Bucket = qiniuyunConfig["bucket"]
	qiniuyunoss.SecretKey = qiniuyunConfig["secretKey"]
	qiniuyunoss.AccessKey = qiniuyunConfig["accessKey"]
	qiniuyunoss.Url = qiniuyunConfig["url"]
	qiniuyunoss.Load()

	sentinel.Rules = globalConfig.GetStringMap("SentinelRules")
	sentinel.Load()

	service.TempVideoFolderPath = globalConfig.GetString("TempVideoFolderPath")
}
