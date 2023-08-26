package pg

import (
	"context"
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
		newVideo, _ := video.ToApiVideo(DB, clientUser, false)
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

	saveCoverPath, dbCoverPath := utils.GetVideoCoverName(saveName)

	video := DBVideo{Author: user.ID, Title: request.Title, PlayUrl: dbURL, CoverUrl: dbCoverPath}

	tx := DB.Begin()
	res := video.insert(tx)

	if !res {
		response.StatusCode = 2
		str := utils.ErrDBSaveVideoFaile.Error()
		response.StatusMsg = &str
		tx.Rollback()
		return
	}

	res = user.increaseWork(tx, 1)

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

	go utils.ExtractCover(saveName, saveCoverPath)
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
		newVideo, _ := video.ToApiVideo(DB, clientUser, false)
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

// 喜欢操作
func DBFavoriteAction(request *api.FavoriteActionRequest, response *api.FavoriteActionResponse, ctx context.Context) {
	// 校验 VideoID
	if request.VideoID <= 0 {
		response.StatusCode = 2
		str := utils.ErrWrongParam.Error()
		response.StatusMsg = &str
		return
	}
	// 校验 token
	_, err := ValidateToken(request.Token)
	if err != nil {
		response.StatusCode = 1
		str := utils.ErrTokenVerifiedFailed.Error()
		response.StatusMsg = &str
		return
	}
	clientUser, _ := ValidateToken(request.Token)
	newRecord := &Like{
		VideoId: request.VideoID,
		UserId:  clientUser.ID,
	}

	curVideo := &DBVideo{
		ID: request.VideoID,
	}

	var ans bool // 点赞或取消点赞是否成功

	// 更改favorited_videos表时，需要同步更改users表盒videos表
	if request.ActionType == 1 {
		// 点赞
		ans = newRecord.insert(DB, clientUser, curVideo)

	} else if request.ActionType == 2 {
		// 取消点赞
		ans = newRecord.delete(DB, clientUser, curVideo)
	} else {
		ans = false
	}
	if !ans {
		response.StatusCode = 2
		str := utils.ErrLikeFaile.Error()
		response.StatusMsg = &str
		return
	}

	response.StatusCode = 0
	response.StatusMsg = &utils.FavoriteVideoActionSuccess

}

// 获取喜欢列表
func DBFavoriteList(request *api.FavoriteListRequest, response *api.FavoriteListResponse) {

	// 如果UserID不合法，直接返回
	if request.UserID <= 0 {
		response.VideoList = nil
		response.StatusCode = 2
		str := utils.ErrWrongParam.Error()
		response.StatusMsg = &str
		return
	}

	// 验证token
	clientUser, err := ValidateToken(request.Token)
	if err != nil || clientUser.ID != request.UserID {
		response.StatusCode = 1
		str := utils.ErrWrongToken.Error()
		response.StatusMsg = &str
		return
	}
	likeobj := &Like{UserId: request.UserID}

	// 根据映射关系查询喜欢列表
	dbvlist, find := likeobj.QueryVideoByUser(DB)
	if !find {
		// 没有找到
		response.StatusCode = 2
		str := utils.ErrGetVideoFromUSer.Error()
		response.StatusMsg = &str
		return
	}

	for _, item := range dbvlist {
		// true 表示确信当前的视频被喜欢
		apiVideo, _ := item.ToApiVideo(DB, clientUser, true)
		response.VideoList = append(response.VideoList, apiVideo)
	}
	response.StatusCode = 0
	response.StatusMsg = &utils.FavoriteVideoListSuccess
}

// 添加评论和删除评论
func DBCommentAction(request *api.CommentActionRequest, response *api.CommentActionResponse, ctx context.Context) {
	// 如果VideoID不合法，直接返回
	if request.VideoID <= 0 {
		response.Comment = nil
		response.StatusCode = 2
		str := utils.ErrWrongParam.Error()
		response.StatusMsg = &str
		return
	}

	// 验证 token
	clientUser, err := ValidateToken(request.Token)
	if err != nil {
		return
	}

	// 发布评论
	if request.ActionType == 1 {
		// 判断评论内容是否是否为空，为空则直接返回
		if len(*request.CommentText) == 0 {
			response.StatusCode = 2
			str := utils.ErrWrongParam.Error()
			response.StatusMsg = &str
			return
		}
		comment, _ := dbCreateComment(request, clientUser.ID)
		response.StatusCode = 0
		cUser := &DBUser{ID: comment.UserId}
		response.Comment, err = comment.ToApiComment(cUser, clientUser) // 删除评论
		// 删除评论
	} else if request.ActionType == 2 {
		// 校验CommentID 是否合法
		if *request.CommentID <= 0 {
			response.StatusCode = 2
			str := utils.ErrWrongParam.Error()
			response.StatusMsg = &str
			return
		}
		if err = DeleteComment(*request.CommentID); err != nil {
			return
		}
		response.StatusCode = 0
		response.StatusMsg = &utils.DeleteCommentSuccess
	} else {
		err = utils.ErrTypeNotSupport
	}

	// 评论计数
	dbv := &DBVideo{ID: request.VideoID}
	txn := DB.Begin()
	dbv.increaseComment(txn)
	txn.Commit()
}

// 获取评论列表
func DBCommentList(request *api.CommentListRequest, response *api.CommentListResponse) {
	if request.VideoID <= 0 {
		response.StatusCode = 2
		str := utils.ErrWrongParam.Error()
		response.StatusMsg = &str
		return
	}

	comments, _ := NewGetCommentListService(request.VideoID, request.Token)
	response.CommentList = comments
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
