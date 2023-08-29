package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type DailyResetHook struct {
	logFile   *os.File
	resetHour int
}

func NewDailyResetHook(logFile *os.File, resetHour int) *DailyResetHook {
	return &DailyResetHook{
		logFile:   logFile,
		resetHour: resetHour,
	}
}

func (hook *DailyResetHook) Fire(entry *logrus.Entry) error {
	now := time.Now()
	if now.Hour() == hook.resetHour {
		// 关闭之前的日志文件
		hook.logFile.Close()

		// 打开新的日志文件
		// 获取当前日期作为日志输出文件
		currentTime := time.Now()
		currentDate := currentTime.Format("2006-01-02")
		fileName := currentDate + ".log"
		newLogFile, err := os.OpenFile(filepath.Join("./log", fileName), os.O_APPEND|os.O_WRONLY|os.O_EXCL, 0666)
		if err != nil {
			return err
		}
		hook.logFile = newLogFile

		// 更新 Logrus 的输出
		entry.Logger.Out = newLogFile
	}
	return nil
}

func (hook *DailyResetHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
