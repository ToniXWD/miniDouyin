package log

import (
	log "github.com/sirupsen/logrus"
	"miniDouyin/utils"
	"os"
	"path/filepath"
	"time"
)

func Init() {
	log.SetLevel(utils.LOG_LEVEL)          // 日志级别
	log.SetFormatter(&log.JSONFormatter{}) // 日志格式

	if utils.USE_TERMINAL {
		log.SetOutput(os.Stdout)
	} else {
		// 获取当前日期作为日志输出文件
		currentTime := time.Now()
		currentDate := currentTime.Format("2006-01-02")
		fileName := currentDate + ".log"
		file, err := os.OpenFile(filepath.Join("./log", fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic("初始化日志记录失败！")
		}
		log.SetOutput(file)

		if utils.UPDATE_DAILY {
			// 每天要更新日志文件
			// 创建每天定时任务的 Hook
			resetHook := NewDailyResetHook(file, 0) // 在每天的 0 点重新记录
			log.AddHook(resetHook)
		}
	}
}
