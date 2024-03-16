package db

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"

	"work/pkg/constants"
)

var DB *gorm.DB

func Load() {
	var err error
	DB, err = gorm.Open(mysql.Open(constants.MysqlDSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	if err = DB.Use(gormopentracing.New()); err != nil {
		panic(err)
	}

	hlog.Infof("Mysql connected successfully.")
}
