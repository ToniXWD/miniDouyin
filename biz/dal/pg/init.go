package pg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"miniDouyin/utils"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(postgres.Open(utils.DSN),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		panic(err)
	}
	ChanFromDB = make(chan RedisMsg)
	go PGserver()
}
