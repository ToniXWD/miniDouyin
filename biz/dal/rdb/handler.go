package rdb

import (
	"context"
	"miniDouyin/biz/model/miniDouyin/api"
	"strconv"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

// 缓存完成路由业务Login
func RedisLogin(request *api.UserLoginRequest, response *api.UserLoginResponse) bool {
	log.Debugln("Trying to get User from redis!")

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
		isFollow, err := IsFollow(request.UserID, request.UserID)
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

// 从缓存完成业务PublishList
func RedisPublishList(request *api.PublishListRequest, response *api.PublishListResponse) bool {
	ids, valid := GetPublishListByUserID(request.UserID)
	if !valid {
		// 缓存不能处理
		return false
	}

	for _, item := range ids {
		vMap, valid := GetVideoById(item)
		if !valid {
			// 缓存不能处理
			response.VideoList = nil
			return false
		}
		apiV, valid := VMap2ApiVidio(request.UserID, vMap)
		if !valid {
			// 缓存不能处理
			response.VideoList = nil
			return false
		}
		response.VideoList = append(response.VideoList, apiV)
	}

	response.StatusCode = 0
	return true
}

// 缓存完成路由业务Feed
// func RedisFeed(request *api.FeedRequest, response *api.FeedResponse) bool {
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
// }

// 缓存完成 GetCommentList
func RedisGetCommentList(request *api.CommentListRequest, response *api.CommentListResponse) bool {
	// 通过缓存查找评论列表

	// 获取评论列表
	clist, find := GetVideoCommentList(int(request.VideoID))
	if !find {
		return false
	}

	for _, cid := range clist {
		cMap, find := GetCommentByID(cid)
		if !find {
			response.CommentList = nil
			return false
		}
		apiC, conv := CMap2ApiComment(cMap)
		if !conv {
			response.CommentList = nil
			return false
		}
		response.CommentList = append(response.CommentList, apiC)
	}

	response.StatusCode = 0
	str := "Get follow list successfully"
	response.StatusMsg = &str
	return true
}

// 缓存完成 FavoriteList
func RedisGetFavoriteList(request *api.FavoriteListRequest, response *api.FavoriteListResponse) bool {
	// 获取点赞 ID 列表
	ids, valid := GetFavoriteListByUserID(request.UserID)
	if !valid {
		// 缓存不能处理
		return false
	}

	for _, item := range ids {
		// 通过点赞ID获取点赞(cMap)
		lMap, valid := GetLikeByID(item)
		if !valid {
			// 缓存不能处理
			return false
		}
		videoId := lMap["VideoId"]
		vMap, valid := GetVideoById(videoId)
		if !valid {
			// 缓存不能处理
			response.VideoList = nil
			return false
		}
		apiV, valid := VMap2ApiVidio(request.UserID, vMap)
		if !valid {
			// 缓存不能处理
			response.VideoList = nil
			return false
		}
		response.VideoList = append(response.VideoList, apiV)
	}
	response.StatusCode = 0
	return true
}

// 缓存完成关注列表
func RedisGetFollowList(request *api.RelationFollowListRequest, response *api.RelationFollowListResponse) bool {
	// 获取关注列表
	idlist, find := GetFollowsIDList(request.UserID)
	if !find {
		return false
	}

	for _, id := range idlist {
		ID, _ := strconv.Atoi(id)
		tMap, find := GetUserById(int64(ID))
		if !find {
			response.UserList = nil
			return false
		}
		apiU := GetApiUserFromMap(tMap)
		// 自己一定关注了该用户
		apiU.IsFollow = true
		response.UserList = append(response.UserList, apiU)
	}
	response.StatusCode = 0
	return true
}

// 缓存完成粉丝列表
func RedisGetFollowerList(request *api.RelationFollowerListRequest, response *api.RelationFollowerListResponse) bool {
	// 获取粉丝列表
	idlist, find := GetFollowersIDList(request.UserID)
	if !find {
		return false
	}

	for _, id := range idlist {
		ID, _ := strconv.Atoi(id)
		tMap, find := GetUserById(int64(ID))
		if !find {
			response.UserList = nil
			return false
		}
		apiU := GetApiUserFromMap(tMap)
		// 判断自己是否也关注了该粉丝
		isfollow, err := IsFollow(request.UserID, apiU.ID)
		if err != nil {
			response.UserList = nil
			return false
		}
		apiU.IsFollow = isfollow
		response.UserList = append(response.UserList, apiU)
	}
	response.StatusCode = 0
	return true
}

func RedisGetFriendList(request *api.RelationFriendListRequest, response *api.RelationFriendListResponse) bool {
	// 通过缓存查找好友列表
	flist, find := GetFriendList(request.UserID)
	if !find || len(flist) == 0 {
		return false
	}

	for _, fid := range flist {
		fMap, find := GetFriendByID(fid, strconv.Itoa(int(request.UserID)))
		if !find {
			response.UserList = nil
			return false
		}

		user_id, _ := strconv.ParseInt(fid, 10, 64)
		user, find := GetUserById(user_id)
		if !find {
			response.UserList = nil
			return false
		}

		apiF := FMap2ApiFriend(fMap, user)
		response.UserList = append(response.UserList, apiF)
	}

	response.StatusCode = 0
	str := "Get friend list successfully"
	response.StatusMsg = &str
	return true
}

func RedisGetChatRec(request *api.ChatRecordRequest, response *api.ChatRecordResponse) bool {
	ctx := context.Background()
	// 通过缓存查找聊天记录
	log.Debugln("传入的时间戳为 = ", request.PreMsgTime)
	user, _ := GetUserByToken(request.Token)
	cmp := float64(request.PreMsgTime + 1)
	fid := user["ID"]
	tid := strconv.Itoa(int(request.ToUserID))

	if v, _ := strconv.ParseInt(fid, 10, 64); v > request.ToUserID {
		fid, tid = tid, fid
	}

	chatRec_key := "chatrec_" + fid + "_" + tid
	res, err := Rdb.ZRangeByScore(ctx, chatRec_key, &redis.ZRangeBy{
		Min:    strconv.Itoa(int(cmp)),
		Max:    "+inf",
		Offset: 0,
		Count:  -1,
	}).Result()
	if err != nil {
		response.MessageList = nil
		log.Debugln(err.Error())
		return false
	}

	// fromUserID, _ := strconv.ParseInt(user["ID"], 10, 64)
	// chatList, _, find := GetChatRec(fromUserID, request.ToUserID)
	// if !find {
	// 	return false
	// }

	for _, msg_id := range res {
		content_key := "content_" + msg_id
		content, _ := Rdb.HGetAll(ctx, content_key).Result()
		apiC := CMap2ApiChat(content)

		response.MessageList = append(response.MessageList, apiC)
	}

	response.StatusCode = 0
	str := "Get chat record successfully"
	response.StatusMsg = &str
	return true
}
