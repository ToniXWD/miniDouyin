/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-09-02 20:38:32
 * @LastEditTime: 2023-09-03 15:36:11
 * @version: 1.0
 */
package rdb

import (
	"context"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func UpdateChatRec(data map[string]interface{}) {
	ctx := context.Background()
	// 设置key
	fid, _ := data["FromID"].(int64)
	tid, _ := data["ToID"].(int64)
	createAt, _ := data["CreatedAt"].(int64)
	if fid > tid {
		fid, tid = tid, fid
	}
	chatrec_key := "chatrec_" + strconv.Itoa(int(fid)) + "_" + strconv.Itoa(int(tid))

	item := redis.Z{
		Member: data["ID"].(int64),
		Score:  float64(createAt),
	}
	_, err := Rdb.ZAdd(ctx, chatrec_key, item).Result()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, chatrec_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))

	content := map[string]interface{}{
		"ID":        data["ID"].(int64),
		"FromID":    data["FromID"].(int64),
		"ToID":      data["ToID"].(int64),
		"CreatedAt": createAt,
		"Content":   data["Content"].(string),
	}
	// 设置key
	content_key := "content_" + strconv.Itoa(int(data["ID"].(int64)))
	_, err = Rdb.HMSet(ctx, content_key, content).Result()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, content_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))

}

func GetChatRec(fid int64, tid int64) ([]string, string, bool) {
	ctx := context.Background()
	if fid > tid {
		fid, tid = tid, fid
	}
	chatrec_key := "chatrec_" + strconv.Itoa(int(fid)) + "_" + strconv.Itoa(int(tid))

	// 使用 Exists 方法判断键是否存在
	_, err := Rdb.Exists(ctx, chatrec_key).Result()
	if err != nil {
		log.Debugln("Error:", err)
		return nil, "", false
	}

	chatrec, err := Rdb.ZRange(ctx, chatrec_key, 0, -1).Result()
	if err != nil || len(chatrec) == 0 {
		return nil, "", false
	}
	return chatrec, chatrec_key, true
}

func CMap2ApiChat(cMap map[string]string) (apiChat *api.Message) {
	id, _ := strconv.ParseInt(cMap["ID"], 10, 64)
	toUserID, _ := strconv.ParseInt(cMap["ToID"], 10, 64)
	fromUserID, _ := strconv.ParseInt(cMap["FromID"], 10, 64)
	createTime, _ := strconv.ParseInt(cMap["CreatedAt"], 10, 64)
	apiChat = &api.Message{
		ID:         id,
		ToUserID:   toUserID,
		FromUserID: fromUserID,
		Content:    cMap["Content"],
		CreateTime: &createTime,
	}
	return apiChat

}
