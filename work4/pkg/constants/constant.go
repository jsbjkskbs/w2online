package constants

import (
	"time"
)

// 数据库连接
const (
	MysqlDSN = `work:work123456@tcp(localhost:13306)/work`

	RedisAddr     = `localhost:16379`
	RedisPassword = `work123456`
)

const (
	DefaultPageSize        = 10
	ESNoKeywordsFlag       = ``
	ESNoTimeFilterFlag     = -1
	ESNoUsernameFilterFlag = ``
	ESNoPageParamFlag      = -1
)

// 默认url
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
