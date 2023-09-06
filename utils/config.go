/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-08-23 11:10:17
 * @LastEditTime: 2023-08-23 11:10:50
 * @version: 1.0
 */
package utils

import  (
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	// 关系型数据库配置项
	DSN    = "user=postgres password="+os.Getenv("POSTGRE_SQL_PASSWORD")+" dbname=postgres host="+os.Getenv("POSTGRE_SQL_HOST")+" port="+os.Getenv("POSTGRE_SQL_PORT")+" sslmode=disable TimeZone=Asia/Shanghai"
	DBTYPE = "pg"

	// hertz配置项
	// PORT     = "8080"
	URLIP    = os.Getenv("paas_url")
	ServerIP = "0.0.0.0"
	MaxBody  = 128 * 1024 * 1024

	// Redis配置项
	REDIS_ADDR     = os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	REDIS_PASSWD   = os.Getenv("REDISCLI_AUTH")
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
