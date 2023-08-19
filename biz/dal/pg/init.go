package pg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dsn = "user=postgres password=tmdgnnwscjl dbname=douyin port=5432 sslmode=disable TimeZone=Asia/Shanghai"

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		panic(err)
	}
}
