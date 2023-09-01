package rdb

import (
	"miniDouyin/utils"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     utils.REDIS_ADDR,
		Password: utils.REDIS_PASSWD, // 没有密码，默认值
		DB:       utils.REDIS_DB,     // 默认DB 0
	})
}
