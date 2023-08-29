package rdb

import (
	"fmt"
	"miniDouyin/biz/model/miniDouyin/api"
	"strconv"
)

// 缓存完成路由业务Login
func RedisLogin(request *api.UserLoginRequest, response *api.UserLoginResponse) bool {
	fmt.Println("Trying to get User from redis!")

	// 数据库中的key是token
	user_key := request.Username + request.Password

	var user map[string]string
	var find bool

	if user, find = GetUserByToken(user_key); !find {
		return false // 缓存不能处理
	}

	// 该key存在则进行校验
	if request.Password != user["Passwd"] {
		response.StatusCode = 1
		str := "Password wrong!"
		response.StatusMsg = &str
		return true
	}
	response.StatusCode = 0
	num, _ := strconv.Atoi((user["ID"]))
	response.UserID = int64(num)
	response.Token = user["Token"]
	str := "Login successfully!"
	response.StatusMsg = &str
	return true
}

// 缓存完成路由业务GetUserInfo
func RedisGetUserInfo(request *api.UserRequest, response *api.UserResponse) bool {
	// 通过缓存查询被查看的用户
	user, find := GetUserById(request.UserID)
	// 如果被查询的用户不存在，直接返回
	if !find {
		return false
	}

	// user是一个map
	response.User = GetApiUserFromMap(user)

	// 查看2者关系是否存在关注关系
	if user["Token"] == request.Token {
		// 如果要查看的用户就是自己，那就是已关注状态
		response.User.IsFollow = true
	} else {
		// 否则要查看关注关系
		isFollow, err := IsFollow(request.Token, request.UserID)
		if err != nil {
			// 缓存中没有token用户的记录，该业务无法通过缓存完成
			return false
		}
		if isFollow {
			//	client关注了指定ID的user
			response.User.IsFollow = true
		} else {
			response.User.IsFollow = false
		}
	}
	response.StatusCode = 0
	str := "Get user information successfully!"
	response.StatusMsg = &str
	return true
}

// 缓存完成路由业务Feed
//func RedisFeed(request *api.FeedRequest, response *api.FeedResponse) bool {
//	var clientUser map[string]string = nil
//	var find bool
//	if request.Token != nil {
//		clientUser, find = GetUserByToken(request.GetToken())
//		if !find {
//			// 没有找到对应的client，直接返回
//			return false
//		}
//	}
//
//	vids, err := GetFeedIds(*request.LatestTime)
//	if err != nil {
//		return false
//	}
//
//	for _, v_id := range vids {
//		video, res := GetVideoById(v_id)
//		if !res {
//			// 没招到就重新从数据库加载
//			return false
//		}
//		apiv, find := VMap2ApiVidio(clientUser, video)
//		if !find {
//			return false
//		}
//		response.VideoList = append(response.VideoList, apiv)
//	}
//	response.StatusCode = 0
//	str := "Load video list successfully"
//	response.StatusMsg = &str
//	return true
//}
