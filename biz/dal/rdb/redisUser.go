package rdb

import (
	"context"
	"fmt"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
	"strconv"
	"time"
)

// 新建User缓存项
func NewUser(data map[string]interface{}) {
	// 同时新建以ID和token为key的项
	ctx := context.Background()
	// 设置key
	token := data["Token"].(string)
	_, err := Rdb.HMSet(ctx, token, data).Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, token, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))

	// 设置key
	ID := strconv.Itoa(int(data["ID"].(int64)))
	_, err = Rdb.HMSet(ctx, ID, data).Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	// 设置过期时间
	Rdb.Expire(ctx, ID, time.Hour*time.Duration(utils.REDIS_HOUR_TTL))
}

// 通过token获取缓存中的记录
func GetUserByToken(token string) (map[string]string, bool) {
	ctx := context.Background()

	// 使用 Exists 方法判断键是否存在
	exists, err := Rdb.Exists(ctx, token).Result()
	if err != nil || exists != 1 {
		fmt.Println("Error:", err)
		return nil, false
	}

	user, err := Rdb.HGetAll(ctx, token).Result()
	if err != nil {
		return user, true
	}
	return nil, false
}

// 通过id获取缓存中的记录
func GetUserById(ID int64) (map[string]string, bool) {
	ctx := context.Background()
	id := strconv.Itoa(int(ID))

	// 使用 Exists 方法判断键是否存在
	exists, err := Rdb.Exists(ctx, id).Result()
	if err != nil || exists != 1 {
		fmt.Println("Error:", err)
		return nil, false
	}

	user, err := Rdb.HGetAll(ctx, id).Result()
	if err != nil {
		return nil, false
	}
	return user, true
}

// 从Map构造apiUser，但IsFollow需要自己另行判断
func GetApiUserFromMap(user map[string]string) *api.User {
	// 填充结构体
	ID, _ := strconv.Atoi(user["ID"])
	tmp, _ := strconv.Atoi(user["FollowCount"])
	FollowCount := int64(tmp)
	tmp, _ = strconv.Atoi(user["FollowerCount"])
	FollowerCount := int64(tmp)
	Avatar := user["Avatar"]
	BackgroundImage := user["BackgroundImage"]
	Signature := user["Signature"]
	tmp, _ = strconv.Atoi(user["TotalFavorited"])
	TotalFavorited := int64(tmp)
	tmp, _ = strconv.Atoi(user["WorkCount"])
	WorkCount := int64(tmp)
	tmp, _ = strconv.Atoi(user["FavoriteCount"])
	FavoriteCount := int64(tmp)

	User := &api.User{
		ID:              int64(ID),
		Name:            user["Username"],
		FollowCount:     &FollowCount,
		FollowerCount:   &FollowerCount,
		IsFollow:        false,
		Avatar:          nil,
		BackgroundImage: nil,
		Signature:       nil,
		TotalFavorited:  &TotalFavorited,
		WorkCount:       &WorkCount,
		FavoriteCount:   &FavoriteCount,
	}
	if Avatar != "" {
		User.Avatar = &Avatar
	}
	if BackgroundImage != "" {
		User.BackgroundImage = &BackgroundImage
	}
	if Signature != "" {
		User.Signature = &Signature
	}

	return User
}
