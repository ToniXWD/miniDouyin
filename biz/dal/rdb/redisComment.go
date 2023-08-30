package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
	"strconv"
	"time"
)

// 新建 Comment 缓存项
func NewComment(data map[string]interface{}) {
	ctx := context.Background()

	creat_at := data["CreatedAt"]
	value, ok := creat_at.(string)
	if ok {
		fmt.Println("Value is a time.Time:", value)
	} else {
		fmt.Println("Value is not a time.Time")
	}
	score, _ := utils.StringToFloat64(value)
	z := redis.Z{
		Score:  score, // CreateAt
		Member: data,
	}

	// 设置key
	vID := strconv.Itoa(int(data["VideoId"].(int64)))
	comment_key := "video_" + "comments" + ":" + vID
	_, err := Rdb.ZAdd(ctx, comment_key, z).Result()
	if err != nil {
		log.Debugln(err.Error())
	}

	// 设置过期时间
	Rdb.Expire(ctx, vID, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 从缓存中删除评论
func DelComment(data map[string]interface{}) {
	ctx := context.Background()
	// 设置key
	vID := strconv.Itoa(int(data["VideoId"].(int64)))
	comment_key := "video_" + "comments" + ":" + vID

	_, err := Rdb.ZRem(ctx, comment_key, data).Result()
	if err != nil {
		log.Debugln(err.Error())
	}
}

// 通过 videoId 获取评论列表
func ListComment(videoId int64) ([]map[string]interface{}, bool) {
	ctx := context.Background()

	// 设置 key
	videoID := strconv.Itoa(int(videoId))
	comment_key := "video_" + "comments" + ":" + videoID
	var commentList []map[string]interface{}
	err := Rdb.ZRevRange(ctx, comment_key, 0, -1).ScanSlice(&commentList)
	if err != nil {
		return nil, false
	}
	return commentList, true
}

// 从 CMap 转化为api.Comment
func CMap2ApiComment(cMap map[string]interface{}) (*api.Comment, bool) {
	ID := cMap["ID"].(int64)
	content := cMap["Content"].(string)
	createdAt := cMap["CreatedAt"].(string)
	userId := cMap["UserId"].(int64)

	// 查询评论用户缓存
	cUserMap, find := GetUserById(userId)
	if !find {
		return nil, false
	}
	cUser := GetApiUserFromMap(cUserMap)

	apic := &api.Comment{
		ID:         ID,
		User:       cUser,
		Content:    content,
		CreateDate: createdAt,
	}
	return apic, true
}
