/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-08-23 11:10:17
 * @LastEditTime: 2023-08-23 11:10:50
 * @version: 1.0
 */
package utils

import log "github.com/sirupsen/logrus"

var (
	// 关系型数据库配置项
	DSN    = "user=toni password=tmdgnnwscjl dbname=douyin port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	DBTYPE = "pg"

	// hertz配置项
	PORT     = "8889"
	URLIP    = "172.17.44.93"
	ServerIP = "0.0.0.0"
	MaxBody  = 128 * 1024 * 1024

	// Redis配置项
	REDIS_ADDR     = "localhost:6379"
	REDIS_PASSWD   = ""
	REDIS_DB       = 0
	REDIS_HOUR_TTL = 5
	REDIS_MAX_FEED = 1000

	// 日志配置项
	USE_TERMINAL = true // 使用终端作为日志输出？
	UPDATE_DAILY = true // 是否每天更新日志文件？
	LOG_LEVEL    = log.DebugLevel

	// Gorm 配置
	GORM_LOGGER_TERMINAL = false // 日志输出到终端？
)
