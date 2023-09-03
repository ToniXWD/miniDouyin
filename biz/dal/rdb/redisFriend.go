/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-09-02 15:22:38
 * @LastEditTime: 2023-09-03 13:50:05
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

func NewFriend(data map[string]interface{}) {
	ctx := context.Background()
	fid := strconv.Itoa(int(data["FromID"].(int64)))
	tid := strconv.Itoa(int(data["ToID"].(int64)))
	if fid > tid {
		fid, tid = tid, fid
	}
	friend_key := "friend_" + fid + "_" + tid
	_, err := Rdb.HMSet(ctx, friend_key, data).Result()
	if err != nil {
		log.Debugln(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, friend_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

func UpdateFriendList(data map[string]interface{}) {
	ctx := context.Background()
	user_id := data["ID"].(int64)
	clientuser_id := strconv.Itoa(int(user_id))
	friend_id := data["Friend"].(int64)
	frilist_key := "friendlist_" + clientuser_id

	item := redis.Z{
		Member: strconv.Itoa(int(friend_id)),
		Score:  float64(friend_id),
	}
	_, err := Rdb.ZAdd(ctx, frilist_key, item).Result()
	if err != nil {
		log.Debugln(err.Error())
	}

	// 设置过期时间
	Rdb.Expire(ctx, frilist_key, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

func GetFriendList(user_id int64) ([]string, bool) {
	ctx := context.Background()
	clientuser_id := strconv.Itoa(int(user_id))
	frilist_key := "friendlist_" + clientuser_id

	// 使用 Exists 方法判断键是否存在
	_, err := Rdb.Exists(ctx, frilist_key).Result()
	if err != nil {
		log.Debugln("Error:", err)
		return nil, false
	}

	frilist, err := Rdb.ZRange(ctx, frilist_key, 0, -1).Result()
	if err != nil || len(frilist) == 0 {
		return nil, false
	}
	return frilist, true
}

func GetFriendByID(fid string, tid string) (map[string]string, bool) {
	ctx := context.Background()
	if fid > tid {
		fid, tid = tid, fid
	}
	friend_key := "friend_" + fid + "_" + tid

	// 判断键是否存在
	_, err := Rdb.Exists(ctx, friend_key).Result()
	if err != nil {
		log.Debugln(err.Error())
		return nil, false
	}

	friend, err := Rdb.HGetAll(ctx, friend_key).Result()
	if err != nil {
		log.Debugln(err.Error())
		return nil, false
	}
	return friend, true
}

func FMap2ApiFriend(fmap map[string]string, user map[string]string) *api.FriendUser {
	apiuser := GetApiUserFromMap(user)
	msg := fmap["Message"]
	msgType, _ := strconv.ParseInt(fmap["MessageType"], 10, 64)
	apif := &api.FriendUser{
		ID:              apiuser.ID,
		Name:            apiuser.Name,
		FollowCount:     apiuser.FollowCount,
		FollowerCount:   apiuser.FollowerCount,
		IsFollow:        apiuser.IsFollow,
		Avatar:          apiuser.Avatar,
		BackgroundImage: apiuser.BackgroundImage,
		Signature:       apiuser.Signature,
		TotalFavorited:  apiuser.TotalFavorited,
		WorkCount:       apiuser.WorkCount,
		FavoriteCount:   apiuser.FavoriteCount,
		Message:         &msg,
		MsgType:         msgType,
	}

	return apif
}
