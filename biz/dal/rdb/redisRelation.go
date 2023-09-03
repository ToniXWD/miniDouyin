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

	UserID := data["UserID"].(int64)
	FollowID := data["FollowID"].(int64)

	// 设置key
	relation_key := "follows_" + strconv.Itoa(int(UserID))
	err := Rdb.SAdd(ctx, relation_key, FollowID).Err()
	if err != nil {
		log.Debugln(err.Error())
	}

	// 设置过期时间
	Rdb.Expire(ctx, relation_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))

	// 设置粉丝关系 key
	follower_key := "followers_" + strconv.Itoa(int(FollowID))
	err = Rdb.SAdd(ctx, follower_key, UserID).Err()
	if err != nil {
		log.Debugln(err.Error())
	}
	Rdb.Expire(ctx, follower_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 删除Follow关系缓存项
func DelFollowRelation(data map[string]interface{}) {
	ctx := context.Background()

	UserID := data["UserID"].(int64)
	FollowID := data["FollowID"].(int64)

	// 设置key
	relation_key := "follows_" + strconv.Itoa(int(UserID))
	err := Rdb.SRem(ctx, relation_key, FollowID).Err()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 删除后重设设置过期时间
	Rdb.Expire(ctx, relation_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))

	// 设置粉丝关系 key
	follower_key := "followers_" + strconv.Itoa(int(FollowID))
	err = Rdb.SRem(ctx, follower_key, UserID).Err()
	if err != nil {
		log.Debugln(err.Error())
	}
	Rdb.Expire(ctx, follower_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 通过缓存判断token用户是否关注了ID用户
func IsFollow(UserID int64, FollowID int64) (bool, error) {
	userID := int(UserID)
	followId := int(FollowID)
	ctx := context.Background()

	relation_key := "follows_" + strconv.Itoa(userID)
	// 使用 Exists 方法判断键是否存在
	exists, err := Rdb.Exists(ctx, relation_key).Result()
	if err != nil || exists != 1 {
		log.Debugln("Error:", err)
		// 关系不再缓存中，则缓存无法处理
		return false, utils.ErrRedisCacheNotFound
	}

	// 判断元素是否在 Set 中
	exi, err := Rdb.SIsMember(ctx, relation_key, followId).Result()
	if err != nil || !exi {
		log.Debugf("缓存查找关注关系成功，%s 没有关注 %s\n", userID, followId)
		return false, nil
	}
	log.Debugf("缓存查找关注关系成功, %s 关注了 %s\n", userID, followId)
	return true, nil
}

// 获取缓存关注列表
func GetFollowsIDList(UserID int64) ([]string, bool) {
	ctx := context.Background()
	follow_key := "follows_" + strconv.Itoa(int(UserID))

	// 使用 Exists 方法判断键是否存在
	exists, err := Rdb.Exists(ctx, follow_key).Result()
	if err != nil || exists != 1 {
		log.Debugln("Error:", err)
		return nil, false
	}
	idList, err := Rdb.SMembers(ctx, follow_key).Result()
	if err != nil {
		return nil, false
	}
	return idList, true
}

// 获取缓存粉丝列表
func GetFollowersIDList(UserID int64) ([]string, bool) {
	ctx := context.Background()
	uid := strconv.Itoa(int(UserID))
	follower_key := "followers_" + uid

	// 使用 Exists 方法判断键是否存在
	exists, err := Rdb.Exists(ctx, follower_key).Result()
	if err != nil || exists != 1 {
		log.Debugln("Error:", err)
		return nil, false
	}
	idlist, err := Rdb.SMembers(ctx, follower_key).Result()
	if err != nil {
		return nil, false
	}
	return idlist, true
}
