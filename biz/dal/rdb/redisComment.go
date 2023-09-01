package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
	"strconv"
	"time"
)

// 新建 Comment 缓存项
func NewComment(data map[string]interface{}) {
	// 同时新建以ID和videoID为key的项
	ctx := context.Background()
	// 设置key
	ID := data["ID"].(int64)
	comment_id := strconv.Itoa(int(ID))
	comment_key := "comment_" + comment_id
	_, err := Rdb.HMSet(ctx, comment_key, data).Result()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, comment_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))

	// 更新comment的id到评论列表
	VID := data["VideoId"].(int64)
	vid := strconv.Itoa(int(VID))
	vid_key := "video_comment_" + vid

	item := redis.Z{
		Member: comment_id,
		Score:  float64(ID),
	}
	_, err = Rdb.ZAdd(ctx, vid_key, item).Result()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, vid_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 通过视频id查询缓存中评论id列表
func GetVideoCommentList(VideoID int) ([]string, bool) {
	ctx := context.Background()
	vid := strconv.Itoa(int(VideoID))
	clist_key := "video_comment_" + vid

	clist, err := Rdb.ZRange(ctx, clist_key, 0, -1).Result()
	if err != nil || len(clist) == 0 {
		return nil, false
	}
	return clist, true
}

// 从缓存中删除评论
func DelComment(data map[string]interface{}) {
	ctx := context.Background()
	// 删除整个comment缓存
	// 设置key
	ID := data["ID"].(int64)
	id := strconv.Itoa(int(ID))
	key := "comment_" + id
	Rdb.Del(ctx, key)

	// 获取视频评论列表的key
	VID := data["VideoId"].(int64)
	vid := strconv.Itoa(int(VID))
	vid_key := "video_comment_" + vid
	// 从视频评论列表里面删除对应的comment的id
	_, err := Rdb.ZRem(ctx, vid_key, data["ID"]).Result()
	if err != nil {
		log.Debugln("删除评论列表中的元素错误", err.Error())
	}

	// 设置过期时间
	Rdb.Expire(ctx, vid_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 通过评论 ID 获取 评论（cMap类型）
func GetCommentByID(id string) (map[string]string, bool) {
	ctx := context.Background()

	key := "comment_" + id

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

// 从 CMap 转化为api.Comment
func CMap2ApiComment(cMap map[string]string) (*api.Comment, bool) {
	strID := cMap["ID"]
	ID, _ := strconv.Atoi(strID)
	content := cMap["Content"]
	createdAt := cMap["CreatedAt"]
	strUserId := cMap["UserId"]
	userId, _ := strconv.Atoi(strUserId)

	// 查询评论用户缓存
	cUserMap, find := GetUserById(int64(userId))
	if !find {
		return nil, false
	}
	cUser := GetApiUserFromMap(cUserMap)

	apic := &api.Comment{
		ID:         int64(ID),
		User:       cUser,
		Content:    content,
		CreateDate: createdAt,
	}
	return apic, true
}
