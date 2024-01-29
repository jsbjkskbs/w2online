package dal

import (
	"work/biz/dal/db"
	"work/biz/mw/redis"
)

func Init() {
	db.Init()
	redis.Init()
}
