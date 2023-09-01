package rdb

import (
	"context"
	log "github.com/sirupsen/logrus"
	"miniDouyin/utils"
	"strconv"
	"time"
)

// 新建Follow关系缓存项
func NewFollowRelation(data map[string]interface{}) {
	ctx := context.Background()
	// 设置key
	relation_key := "follows_" + data["Token"].(string)
	value := data["ID"].(int64)
	err := Rdb.SAdd(ctx, relation_key, value).Err()
	if err != nil {
		log.Debugln(err.Error())
	}

	// 设置过期时间
	Rdb.Expire(ctx, relation_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 删除Follow关系缓存项
func DelFollowRelation(data map[string]interface{}) {
	ctx := context.Background()
	// 设置key
	relation_key := "follows_" + data["Token"].(string)
	err := Rdb.SRem(ctx, relation_key, data["ID"]).Err()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 删除后重设设置过期时间
	Rdb.Expire(ctx, relation_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 通过缓存判断token用户是否关注了ID用户
func IsFollow(token string, ID int64) (bool, error) {
	id := strconv.Itoa(int(ID))
	ctx := context.Background()

	relation_key := "follows_" + token
	// 使用 Exists 方法判断键是否存在
	exists, err := Rdb.Exists(ctx, relation_key).Result()
	if err != nil || exists != 1 {
		log.Debugln("Error:", err)
		// 关系不再缓存中，则缓存无法处理
		return false, utils.ErrRedisCacheNotFound
	}

	// 判断元素是否在 Set 中
	exi, err := Rdb.SIsMember(ctx, relation_key, id).Result()
	if err != nil || !exi {
		log.Debugf("缓存查找关注关系成功，%s 没有关注 %s\n", token, id)
		return false, nil
	}
	log.Debugf("缓存查找关注关系成功, %s 关注了 %s\n", token, id)
	return true, nil
}
