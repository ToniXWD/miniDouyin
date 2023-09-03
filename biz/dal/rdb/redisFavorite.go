package rdb

import (
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"miniDouyin/utils"
	"strconv"
	"time"
)
import "context"

// 新建 Like 缓存项目
func NewLikeVideo(data map[string]interface{}) {

	// 同时新建以ID和userID为key的项
	ctx := context.Background()
	// 设置key
	ID := data["ID"].(int64)
	like_id := strconv.Itoa(int(ID))
	like_key := "like_" + like_id
	_, err := Rdb.HMSet(ctx, like_key, data).Result()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, like_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))

	// UserID:[likeID]
	UID := data["UserId"].(int64)
	uid := strconv.Itoa(int(UID))
	uid_key := "user_like_" + uid

	item := redis.Z{
		Member: like_id,
		Score:  float64(ID),
	}
	_, err = Rdb.ZAdd(ctx, uid_key, item).Result()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, uid_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 从缓存中删除点赞
func DelLikeVideo(data map[string]interface{}) {
	ctx := context.Background()
	// 删除整个like缓存
	// 设置key
	ID := data["ID"].(int64)
	id := strconv.Itoa(int(ID))
	key := "like_" + id
	Rdb.Del(ctx, key)

	// 获取点赞列表的key
	UID := data["uId"].(int64)
	uid := strconv.Itoa(int(UID))
	uid_key := "user_like_" + uid
	// 从视频评论列表里面删除对应的comment的id
	_, err := Rdb.ZRem(ctx, uid_key, data["ID"]).Result()
	if err != nil {
		log.Debugln("删除点赞列表中的元素错误", err.Error())
	}

	// 设置过期时间
	Rdb.Expire(ctx, uid_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 获取缓存点赞ID列表
func GetFavoriteListByUserID(UserID int64) ([]string, bool) {
	ctx := context.Background()
	uid := strconv.Itoa(int(UserID))
	flist_key := "user_like_" + uid

	// 使用 Exists 方法判断键是否存在
	exists, err := Rdb.Exists(ctx, flist_key).Result()
	if err != nil || exists != 1 {
		log.Debugln("Error:", err)
		return nil, false
	}

	clist, err := Rdb.ZRange(ctx, flist_key, 0, -1).Result()
	if err != nil || len(clist) == 0 {
		return nil, false
	}
	return clist, true
}

// 通过点赞ID获取点赞(cMap)
func GetLikeByID(id string) (map[string]string, bool) {
	ctx := context.Background()

	key := "like_" + id

	// 使用 Exists 方法判断键是否存在
	exists, err := Rdb.Exists(ctx, key).Result()
	if err != nil || exists != 1 {
		log.Debugln("Error:", err)
		return nil, false
	}
	comment, err := Rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, false
	}
	return comment, true
}

//func IsVideoLikedById(videoID int64, user_ID int64) (bool, error) {
//	ctx := context.Background()
//
//	user_id := strconv.Itoa(int(ID))
//
//	key := "user_like_" + id
//
//	item := redis.Z{
//		Member: id,
//		Score:  float64(videoID),
//	}
//
//	Rdb.ZCount(ctx, key)
//}
