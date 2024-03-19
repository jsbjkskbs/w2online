package constants

import (
	"time"
)

// 中间件连接
var (
	MysqlDSN                 string
	RabbitmqDSN              string
	RedisAddr, RedisPassword string
	ElasticAddr              string
)

const (
	DefaultPageSize        = 10
	ESNoKeywordsFlag       = ``
	ESNoTimeFilterFlag     = -1
	ESNoUsernameFilterFlag = ``
	ESNoPageParamFlag      = -1
)

// 默认头像url
const (
	DefaultAvatarUrl = ``
)

// 文件大小(以 1 Byte为单位)
const (
	Byte   = 1
	KBytes = 1 * Byte * 1024
	MBytes = 1 * KBytes * 1024
	GBytes = 1 * MBytes * 1024
	TBytes = 1 * GBytes * 1024
	PBytes = 1 * TBytes * 1024
)

const (
	Day  = time.Hour * 24
	Week = Day * 7
)
