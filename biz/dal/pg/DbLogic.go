package pg

import (
	"fmt"
	"mime/multipart"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// 处理登录请求
// 并填充response结构体
func DBUserLogin(request *api.UserLoginRequest, response *api.UserLoginResponse) {
	user := DBUserFromLoginRequest(request)

	if user.QueryUser() {
		// user存在
		fmt.Printf("user = %+v\n", user)

		// 校验密码
		if user.Passwd != request.Password {
			response.StatusCode = 1
			str := "Password wrong!"
			response.StatusMsg = &str
			return
		}
		response.StatusCode = 0
		response.UserID = int64(user.ID)
		response.Token = user.Username + user.Passwd
		str := "Login successfully!"
		response.StatusMsg = &str
		return
	}
	response.StatusCode = 2
	str := "User not exist!"
	response.StatusMsg = &str
}

// 处理注册请求
// 并填充response结构体
func DBUserRegister(request *api.UserRegisterRequest, response *api.UserRegisterResponse) {
	user := DBUserFromRegisterRequest(request)

	if user.QueryUser() {
		// user存在
		response.StatusCode = 1
		str := "User already exists!"
		response.StatusMsg = &str
		return
	}

	if user.insert() {
		response.StatusCode = 0
		response.UserID = int64(user.ID)
		response.Token = user.Token
		str := "Register succesffully!"
		response.StatusMsg = &str
	} else {
		response.StatusCode = 2
		str := "Register failed!"
		response.StatusMsg = &str
	}
}

// 获取User信息
func DBGetUserinfo(request *api.UserRequest, response *api.UserResponse) {
	user, err := DBGetUser(request)
	if err != nil {
		// 没有找到用户或token失败
		response.StatusCode = 1
		str := err.Error()
		response.StatusMsg = &str
		return
	}
	// 填充结构体
	clientUser, _ := ValidateToken(request.Token)
	response.User, _ = user.ToApiUser(clientUser)
	response.StatusCode = 0
	str := "Get user information successfully!"
	response.StatusMsg = &str
}

// 处理视频流
func DBVideoFeed(request *api.FeedRequest, response *api.FeedResponse) {
	vlist, err := GetNewVideoList(*request.LatestTime)
	if err != nil {
		response.StatusCode = 1
		str := utils.ErrGetFeedVideoListFailed.Error()
		response.StatusMsg = &str
		return
	}
	response.StatusCode = 0
	str := "Load video list successfully"
	response.StatusMsg = &str
	for idx, video := range vlist {
		if idx == len(vlist)-1 {
			newNext := utils.TimeToI64(video.CreatedAt)
			response.NextTime = &newNext
		}
		var clientUser *DBUser = nil
		if request.Token != nil {
			clientUser, _ = ValidateToken(*request.Token)
		}
		newVideo, _ := video.ToApiVideo(clientUser)
		response.VideoList = append(response.VideoList, newVideo)
	}
}

// 接受上传视频
func DBReceiveVideo(request *api.PublishActionRequest, response *api.PublishActionResponse, form *multipart.Form, c *app.RequestContext) {
	request.Token = form.Value["token"][0]
	request.Title = form.Value["title"][0]

	// 先验证token
	user, err := ValidateToken(request.Token)

	if err != nil {
		response.StatusCode = 1
		str := utils.ErrTokenVerifiedFailed.Error()
		response.StatusMsg = &str
		return
	}
	file := form.File["data"][0]

	_, saveName, dbURL := utils.GetVideoNameAndPath()

	err = c.SaveUploadedFile(file, saveName)

	if err != nil {
		response.StatusCode = 2
		str := utils.ErrSaveVideoFaile.Error()
		response.StatusMsg = &str
		return
	}

	video := DBVideo{Author: user.ID, Title: request.Title, PlayUrl: dbURL}

	tx := DB.Begin()
	res := video.insert(tx)

	if !res {
		response.StatusCode = 2
		str := utils.ErrDBSaveVideoFaile.Error()
		response.StatusMsg = &str
		tx.Rollback()
		return
	}

	res = user.increaseWork(tx)

	if !res {
		response.StatusCode = 2
		str := utils.ErrDBSaveVideoFaile.Error()
		response.StatusMsg = &str
		tx.Rollback()
		return
	}
	tx.Commit()

	response.StatusCode = 0
	response.StatusMsg = &utils.UploadVideosSuccess
}

// 获取视频播放列表
func DBVideoPublishList(request *api.PublishListRequest, response *api.PublishListResponse) {
	vlist, err := GetUserVideoList(request.UserID)
	if err != nil {
		response.StatusCode = 1
		str := utils.ErrGetUserVideoListFailed.Error()
		response.StatusMsg = &str
		return
	}
	clientUser, _ := ValidateToken(request.Token)
	for _, video := range vlist {
		newVideo, _ := video.ToApiVideo(clientUser)
		response.VideoList = append(response.VideoList, newVideo)
		response.StatusCode = 0
	}
}

// 处理关注请求
// 并填充response结构体
func DBUserAction(request *api.RelationActionRequest, response *api.RelationActionResponse) {
	action := DBActionFromActionRequest(request)

	err := action.ifFollow(request.ActionType)
	fmt.Printf("action = %+v\n", action)
	if err == nil {
		// 关注或取消关注成功
		response.StatusCode = 0
		str := "Action or DeAction successfully!"
		response.StatusMsg = &str
		return
	}
	response.StatusCode = 1
	str := "Action failed!"
	response.StatusMsg = &str
}

func DBFollowList(request *api.RelationFollowListRequest, response *api.RelationFollowListResponse) {
	clientuser, _ := ValidateToken(request.Token)
	var followlist []DBAction
	err := DB.Where("user_id = ?", clientuser.ID).Find(&followlist)
	if err.Error != nil {
		response.StatusCode = 1
		str := "Get follow list failed"
		response.StatusMsg = &str
		response.UserList = nil
		return
	}

	for _, follow := range followlist {
		user := &DBUser{}
		err := DB.Where("id = ?", follow.FollowID).First(user)
		if err.Error != nil {
			response.StatusCode = 1
			str := "Get follow list failed"
			response.StatusMsg = &str
			response.UserList = nil
			return
		}
		apiuser, _ := user.ToApiUser(clientuser)
		response.UserList = append(response.UserList, apiuser)
	}
	response.StatusCode = 0
	str := "Get follow list successfully"
	response.StatusMsg = &str
}

func DBFollowerList(request *api.RelationFollowerListRequest, response *api.RelationFollowerListResponse) {
	clientuser, _ := ValidateToken(request.Token)
	var followerList []DBAction
	err := DB.Where("follow_id = ?", clientuser.ID).Find(&followerList)
	if err.Error != nil {
		response.StatusCode = 1
		str := "Get follower list failed"
		response.StatusMsg = &str
		response.UserList = nil
		return
	}

	for _, follower := range followerList {
		user := &DBUser{}
		err := DB.Where("id = ?", follower.UserID).First(user)
		if err.Error != nil {
			response.StatusCode = 1
			str := "Get follower list failed"
			response.StatusMsg = &str
			response.UserList = nil
			return
		}
		apiuser, _ := user.ToApiUser(clientuser)
		response.UserList = append(response.UserList, apiuser)
	}
	response.StatusCode = 0
	str := "Get follower list successfully"
	response.StatusMsg = &str
}

func DBFriendList(request *api.RelationFriendListRequest, response *api.RelationFriendListResponse) {
	clientuser, _ := ValidateToken(request.Token)
	var friendList []DBfriend
	err := DB.Where("user_id = ?", clientuser.ID).Find(&friendList)
	if err.Error != nil {
		response.StatusCode = 1
		str := "Get friend list failed"
		response.StatusMsg = &str
		response.UserList = nil
		return
	}

	for _, friend := range friendList {
		user := &DBUser{}
		err := DB.Where("id = ?", friend.FriendID).First(user)
		if err.Error != nil {
			response.StatusCode = 1
			str := "Get friend list failed"
			response.StatusMsg = &str
			response.UserList = nil
			return
		}
		apiuser, _ := user.ToApiUser(clientuser)
		response.UserList = append(response.UserList, apiuser)
	}
	response.StatusCode = 0
	str := "Get friend list successfully"
	response.StatusMsg = &str
}

func DBSendMsg(request *api.SendMsgRequest, response *api.SendMsgResponse) {
	if request.ActionType == 1 {
		if sendMsg(request.Token, request.ToUserID, request.Content) {
			response.StatusCode = int32(request.ToUserID)
			response.StatusMsg = 0
		}
	}

	response.StatusMsg = 1
}

func DBChatRec(request *api.ChatRecordRequest, response *api.ChatRecordResponse) {
	clientuser, _ := ValidateToken(request.Token)
	var msgList []DBMessage
	err := DB.Where("from_id = ? AND to_id = ?", clientuser.ID, request.ToUserID).Find(&msgList)
	if err.Error != nil {
		response.StatusCode = 1
		str := "Get chat record failed"
		response.StatusMsg = &str
		response.StructList = nil
		return
	}

	for _, msg := range msgList {
		apimsg := msg.ToApiMessage()
		response.StructList = append(response.StructList, apimsg)
	}
	response.StatusCode = 0
	str := "Get chat record successfully"
	response.StatusMsg = &str
}
