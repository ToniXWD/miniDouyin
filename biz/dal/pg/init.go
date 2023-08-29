package pg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"miniDouyin/utils"
	"os"
	"path/filepath"
)

var DB *gorm.DB

func Init() {
	var err error
	var dbLogger logger.Interface

	if utils.GORM_LOGGER_TERMINAL {
		// 使用终端日志
		dbLogger = logger.Default.LogMode(logger.Info)
	} else {
		// 创建日志文件
		logPath := filepath.Join("./log", "gorm.log")
		file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}
		// 配置日志记录器
		dbLogger = logger.New(
			log.New(file, "\r\n", log.LstdFlags), // 使用 log 包创建新的 logger
			logger.Config{
				SlowThreshold:             200e6, // 慢查询阈值，200ms
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,  // 忽略记录不存在的错误
				Colorful:                  false, // 彩色输出
			},
		)
	}

	DB, err = gorm.Open(postgres.Open(utils.DSN),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 dbLogger,
		})
	if err != nil {
		panic(err)
	}
	ChanFromDB = make(chan RedisMsg)
	go PGserver()
}
