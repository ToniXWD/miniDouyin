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
func NewLike(data map[string]interface{}) {
	// 同时新建以ID和videoID为key的项
	ctx := context.Background()
	// 设置key
	ID := data["ID"].(int64)
	like_id := strconv.Itoa(int(ID))
	like_key := "comment_" + like_id
	_, err := Rdb.HMSet(ctx, like_key, data).Result()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, like_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))

	// 更新like的id到点赞列表
	VID := data["VideoId"].(int64)
	vid := strconv.Itoa(int(VID))
	vid_key := "video_like_" + vid

	item := redis.Z{
		Member: like_id,
		Score:  float64(ID),
	}
	_, err = Rdb.ZAdd(ctx, vid_key, item).Result()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, vid_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 从缓存中删除点赞
func DelLike(data map[string]interface{}) {
	ctx := context.Background()
	// 删除整个like缓存
	// 设置key
	ID := data["ID"].(int64)
	id := strconv.Itoa(int(ID))
	key := "like_" + id
	Rdb.Del(ctx, key)

	// 获取视频点赞列表的key
	VID := data["VideoId"].(int64)
	vid := strconv.Itoa(int(VID))
	vid_key := "video_like_" + vid
	// 从视频点赞列表里面删除对应的like的id
	_, err := Rdb.ZRem(ctx, vid_key, data["ID"]).Result()
	if err != nil {
		log.Debugln("删除点赞列表中的元素错误", err.Error())
	}

	// 设置过期时间
	Rdb.Expire(ctx, vid_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 获取缓存点赞
func GetFavoriteList(VideoID int) ([]string, bool) {
	ctx := context.Background()
	vid := strconv.Itoa(int(VideoID))
	clist_key := "video_like_" + vid

	clist, err := Rdb.ZRange(ctx, clist_key, 0, -1).Result()
	if err != nil || len(clist) == 0 {
		return nil, false
	}
	return clist, true
}

// 通过点赞 ID 获取点赞
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
