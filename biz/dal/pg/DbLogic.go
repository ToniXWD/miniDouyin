/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-09-01 18:42:48
 * @LastEditTime: 2023-09-03 14:49:45
 * @version: 1.0
 */
package pg

import (
	"mime/multipart"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/cloudwego/hertz/pkg/app"
)

// 处理登录请求
// 并填充response结构体
func DBUserLogin(request *api.UserLoginRequest, response *api.UserLoginResponse) {
	user := DBUserFromLoginRequest(request)

	if user.QueryUser() {
		// user存在
		log.Debugf("user = %+v\n", user)

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
		// 发送消息更新缓存
		items := utils.StructToMap(&user)
		msg := RedisMsg{
			TYPE: UserInfo,
			DATA: items,
		}
		ChanFromDB <- msg

		// select {
		// case ChanFromDB <- msg:
		//	log.Debugln("Sent data to ChanFromDB.")
		// default:
		//	log.Debugln("ChanFromDB is not ready for sending.")
		// }

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

		// 发送消息更新缓存
		items := utils.StructToMap(&user)
		msg := RedisMsg{
			TYPE: UserInfo,
			DATA: items,
		}
		ChanFromDB <- msg

	} else {
		response.StatusCode = 2
		str := "Register failed!"
		response.StatusMsg = &str
	}
}

// 获取User信息
func DBGetUserinfo(request *api.UserRequest, response *api.UserResponse) {
	// 先尝试从缓存查找用户
	user, err := ID2DBUser(request.UserID)
	if err != nil {
		response.StatusCode = 3
		str := err.Error()
		response.StatusMsg = &str
		return
	}

	// 填充结构体
	// 先获取client
	// 尝试从缓存获取client
	clientUser, err := Token2DBUser(request.Token)
	if err != nil {
		response.StatusCode = 1
		str := err.Error()
		response.StatusMsg = &str
		return
	}

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
	// 先获取client
	// 尝试从缓存获取client
	clientUser := &DBUser{}
	if request.Token == nil {
		// 未登录状态
		clientUser = nil
	} else {
		// 获取clientUser
		clientUser, err = Token2DBUser(*request.Token)
		if err != nil {
			response.StatusCode = 1
			str := err.Error()
			response.StatusMsg = &str
			return
		}
	}
	for idx, video := range vlist {
		// 遍历视频时，顺手发送消息更新视频缓存
		items := utils.StructToMap(&video)
		msg := RedisMsg{
			TYPE: VideoInfo,
			DATA: items,
		}
		ChanFromDB <- msg

		if idx == len(vlist)-1 {
			newNext := video.CreatedAt.UnixMilli()
			response.NextTime = &newNext
		}
		newVideo, _ := video.ToApiVideo(DB, clientUser, false)
		response.VideoList = append(response.VideoList, newVideo)
	}
	response.StatusCode = 0
	str := "Load video list successfully"
	response.StatusMsg = &str
}

// 接受上传视频
func DBReceiveVideo(request *api.PublishActionRequest, response *api.PublishActionResponse, form *multipart.Form, c *app.RequestContext) {
	request.Token = form.Value["token"][0]
	request.Title = form.Value["title"][0]

	// 先验证token
	// 先尝试缓存获取clientUser
	clientUser, err := Token2DBUser(request.Token)
	if err != nil {
		response.StatusCode = 1
		str := err.Error()
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

	video := DBVideo{Author: clientUser.ID, Title: request.Title, PlayUrl: dbURL, CoverUrl: dbCoverPath}

	tx := DB.Begin()
	res := video.insert(tx)

	if !res {
		response.StatusCode = 2
		str := utils.ErrDBSaveVideoFaile.Error()
		response.StatusMsg = &str
		tx.Rollback()
		return
	}

	res = clientUser.increaseWork(tx, 1)

	if !res {
		response.StatusCode = 2
		str := utils.ErrDBSaveVideoFaile.Error()
		response.StatusMsg = &str
		tx.Rollback()
		return
	}
	tx_res := tx.Commit()
	if tx_res.Error != nil {
		response.StatusCode = 2
		str := utils.ErrDBSaveVideoFaile.Error()
		response.StatusMsg = &str
		tx.Rollback()
		return
	}

	// 事务成功后需要将user的work_count自增
	clientUser.WorkCount++

	// 发送消息更新缓存
	// 由于视频发布后user的字段发送了变化，需要更新缓存的user
	clientUser.UpdateRedis()

	// 更新视频缓存、用户发布视频的集合
	video.UpdateRedis(video.Author)

	response.StatusCode = 0
	response.StatusMsg = &utils.UploadVideosSuccess

	go utils.ExtractCover(saveName, saveCoverPath)
}

// 获取视频发布列表
func DBVideoPublishList(request *api.PublishListRequest, response *api.PublishListResponse) {
	vlist, err := GetUserVideoList(request.UserID)
	if err != nil {
		response.StatusCode = 1
		str := utils.ErrGetUserVideoListFailed.Error()
		response.StatusMsg = &str
		return
	}

	// 先获取client
	// 尝试从缓存获取client
	clientUser, err := Token2DBUser(request.Token)
	if err != nil {
		response.StatusCode = 1
		str := err.Error()
		response.StatusMsg = &str
		return
	}

	// 准备将发布列表更新到缓存
	ids := make([]interface{}, len(vlist))

	for idx, video := range vlist {
		newVideo, _ := video.ToApiVideo(DB, clientUser, false)
		response.VideoList = append(response.VideoList, newVideo)
		response.StatusCode = 0
		ids[idx] = video.ID
		// 更新视频缓存
		items := utils.StructToMap(&video)
		msgV := RedisMsg{
			TYPE: VideoInfo,
			DATA: items,
		}
		ChanFromDB <- msgV
	}
	// 同时更新用户发布的集合缓存
	if clientUser != nil {

		msg := RedisMsg{
			TYPE: Publish,
			DATA: map[string]interface{}{
				"ID":     clientUser.ID,
				"Videos": ids,
			}}
		ChanFromDB <- msg
	}
}

// 处理关注请求
// 并填充response结构体
func DBUserAction(request *api.RelationActionRequest, response *api.RelationActionResponse) {
	action := DBActionFromActionRequest(request)

	err := action.ifFollow(request.ActionType)
	log.Debugf("action = %+v\n", action)
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
func DBFavoriteAction(request *api.FavoriteActionRequest, response *api.FavoriteActionResponse) {
	// 校验 VideoID
	if request.VideoID <= 0 {
		response.StatusCode = 2
		str := utils.ErrWrongParam.Error()
		response.StatusMsg = &str
		return
	}
	// 校验 token
	// 先尝试缓存获取clientUser
	clientUser, err := Token2DBUser(request.Token)
	if err != nil {
		response.StatusCode = 1
		str := err.Error()
		response.StatusMsg = &str
		return
	}

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
		// 发送消息更新缓存
		items := utils.StructToMap(newRecord)
		msg := RedisMsg{
			TYPE: LikeCreate,
			DATA: items,
		}
		ChanFromDB <- msg

	} else if request.ActionType == 2 {
		// 取消点赞
		ans = newRecord.delete(DB, clientUser, curVideo)
		items := utils.StructToMap(newRecord)
		log.Debugf("%T", items["CreateAt"])
		msg := RedisMsg{
			TYPE: LikeDel,
			DATA: items,
		}
		ChanFromDB <- msg
	} else {
		ans = false
	}
	if !ans {
		response.StatusCode = 2
		str := utils.ErrLikeFailed.Error()
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
	// 先尝试缓存获取clientUser
	clientUser, err := Token2DBUser(request.Token)
	if err != nil || clientUser.ID != request.UserID {
		response.StatusCode = 1
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

	for _, video := range dbvlist {
		// 更新视频缓存
		items := utils.StructToMap(&video)
		msgV := RedisMsg{
			TYPE: VideoInfo,
			DATA: items,
		}
		ChanFromDB <- msgV

		// true 表示确信当前的视频被喜欢
		apiVideo, _ := video.ToApiVideo(DB, clientUser, true)
		response.VideoList = append(response.VideoList, apiVideo)
	}
	response.StatusCode = 0
	response.StatusMsg = &utils.FavoriteVideoListSuccess
}

// 添加评论和删除评论
func DBCommentAction(request *api.CommentActionRequest, response *api.CommentActionResponse) {
	// 如果VideoID不合法，直接返回
	if request.VideoID <= 0 {
		response.Comment = nil
		response.StatusCode = 2
		str := utils.ErrWrongParam.Error()
		response.StatusMsg = &str
		return
	}

	// 验证 token
	clientUser, err := Token2DBUser(request.Token)
	if err != nil {
		response.StatusCode = 1
		str := err.Error()
		response.StatusMsg = &str
		return
	}

	if request.ActionType == 1 {
		// 发布评论
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
		response.Comment, err = comment.ToApiComment(cUser, clientUser)
		// 发送消息更新缓存
		items := utils.StructToMap(comment)
		log.Debugf("%T", items["CreateAt"])
		msg := RedisMsg{
			TYPE: CommentCreate,
			DATA: items,
		}
		ChanFromDB <- msg

	} else if request.ActionType == 2 {
		// 删除评论
		// 校验CommentID 是否合法
		if *request.CommentID <= 0 {
			response.StatusCode = 2
			str := utils.ErrWrongParam.Error()
			response.StatusMsg = &str
			return
		}
		// 从数据库删除
		_, err := DeleteComment(*request.CommentID)
		if err != nil {
			return
		}
		// 再从缓存中删除
		msg := RedisMsg{
			TYPE: CommentDel,
			DATA: map[string]interface{}{
				"ID":      *request.CommentID,
				"VideoId": request.VideoID,
			},
		}
		ChanFromDB <- msg

		response.StatusCode = 0
		response.StatusMsg = &utils.DeleteCommentSuccess
	} else {
		err = utils.ErrTypeNotSupport
	}

	// 评论计数
	// 视频是一定存在的
	dbv, _ := GetVideoByID(request.VideoID)

	txn := DB.Begin()
	dbv.increaseComment(txn)
	txn.Commit()

	// 更新视频缓存
	dbv.CommentCount++
	dbv.UpdateRedis(0)
}

// 获取评论列表
func DBCommentList(request *api.CommentListRequest, response *api.CommentListResponse) {
	if request.VideoID <= 0 {
		response.StatusCode = 2
		str := utils.ErrWrongParam.Error()
		response.StatusMsg = &str
		return
	}
	// ToDo:先从缓存取
	comments, _ := NewGetCommentListService(request.VideoID, request.Token)
	response.CommentList = comments
}

// 用户关注列表
func DBFollowList(request *api.RelationFollowListRequest, response *api.RelationFollowListResponse) {
	// 先验证token
	// 先尝试缓存获取clientUser
	clientUser, Terr := Token2DBUser(request.Token)
	if Terr != nil {
		response.StatusCode = 3
		str := Terr.Error()
		response.StatusMsg = &str
		return
	}

	var followlist []DBAction
	err := DB.Where("user_id = ?", clientUser.ID).Find(&followlist)
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
		apiuser, _ := user.ToApiUser(clientUser)

		response.UserList = append(response.UserList, apiuser)
		// 更新缓存
		msg := RedisMsg{
			TYPE: UserFollowAdd,
			DATA: map[string]interface{}{
				"UserID":   follow.UserID, // 粉丝
				"FollowID": follow.FollowID,
			}}
		ChanFromDB <- msg
	}
	response.StatusCode = 0
	str := "Get follow list successfully"
	response.StatusMsg = &str
}

// 用户粉丝列表
func DBFollowerList(request *api.RelationFollowerListRequest, response *api.RelationFollowerListResponse) {
	// 验证Token
	clientUser, Terr := Token2DBUser(request.Token)
	if Terr != nil {
		response.StatusCode = 1
		str := Terr.Error()
		response.StatusMsg = &str
		return
	}
	var followerList []DBAction
	err := DB.Where("follow_id = ?", clientUser.ID).Find(&followerList)
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
		// 更新关注缓存
		msg := RedisMsg{
			TYPE: UserFollowAdd,
			DATA: map[string]interface{}{
				"UserID":   follower.UserID, // 粉丝
				"FollowID": follower.FollowID,
			}}
		ChanFromDB <- msg

		apiuser, _ := user.ToApiUser(clientUser)
		response.UserList = append(response.UserList, apiuser)
	}
	response.StatusCode = 0
	str := "Get follower list successfully"
	response.StatusMsg = &str
}

// 用户好友列表
func DBFriendList(request *api.RelationFriendListRequest, response *api.RelationFriendListResponse) {
	clientUser, Terr := Token2DBUser(request.Token)
	if Terr != nil {
		response.StatusCode = 1
		str := Terr.Error()
		response.StatusMsg = &str
		return
	}
	var friendList []DBfriend
	err := DB.Where("user_id = ?", clientUser.ID).Find(&friendList)
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
		apiuser, _ := user.ToApiUser(clientUser)
		apiFriend, content := apiUser2apiFriend(apiuser, clientUser)
		response.UserList = append(response.UserList, apiFriend)

		item := map[string]interface{}{
			"FromID":  clientUser.ID,
			"ToID":    friend.FriendID,
			"Message": content.Content,
			"MsgType": apiFriend.MsgType,
		}
		msg := RedisMsg{
			TYPE: Friend,
			DATA: item,
		}
		ChanFromDB <- msg

		item = map[string]interface{}{
			"ID":     clientUser.ID,
			"Friend": friend.FriendID,
		}
		msg = RedisMsg{
			TYPE: FriendList,
			DATA: item,
		}
		ChanFromDB <- msg

	}
	response.StatusCode = 0
	str := "Get friend list successfully"
	response.StatusMsg = &str
}

func DBSendMsg(request *api.SendMsgRequest, response *api.SendMsgResponse) {
	if request.ActionType == 1 {
		if content, send := sendMsg(request.Token, request.ToUserID, request.Content); send {

			item := map[string]interface{}{
				"ID":        content.ID,
				"FromID":    content.FromID,
				"ToID":      content.ToID,
				"Content":   content.Content,
				"CreatedAt": content.CreatedAt.UnixMilli(),
			}
			msg := RedisMsg{
				TYPE: ChatRecord,
				DATA: item,
			}
			ChanFromDB <- msg

			item = map[string]interface{}{
				"FromID":  content.FromID,
				"ToID":    content.ToID,
				"Message": content.Content,
				"MsgType": 1,
			}
			msg = RedisMsg{
				TYPE: Friend,
				DATA: item,
			}
			ChanFromDB <- msg
			response.StatusCode = 0
			str := utils.SendMessageSuccess
			response.StatusMsg = &str
			return
		}
	}

	response.StatusCode = 1
}

func DBChatRec(request *api.ChatRecordRequest, response *api.ChatRecordResponse) {
	clientUser, Terr := Token2DBUser(request.Token)
	if Terr != nil {
		response.StatusCode = 1
		str := Terr.Error()
		response.StatusMsg = &str
		return
	}
	var msgList []DBMessage
	log.Debugln("传入的时间戳为 = ", request.PreMsgTime)

	// cmp := time.Unix(request.PreMsgTime, 0)
	var cmp time.Time
	cmp = time.UnixMilli(0)

	if request.PreMsgTime != 0 {
		cmp = time.UnixMilli(request.PreMsgTime)
		cmp = cmp.Add(time.Second * 2)
	}

	log.Debugf("DBChatRec cmp = %v\n", cmp)

	err := DB.Where("((from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)) AND created_at > ?", request.ToUserID, clientUser.ID, clientUser.ID, request.ToUserID, cmp).Order("ID").Find(&msgList)
	// err := DB.Where("((from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?))", request.ToUserID, clientuser.ID, clientuser.ID, request.ToUserID).Order("ID desc").Find(&msgList)
	if err.Error != nil {
		response.StatusCode = 1
		str := "Get chat record failed"
		response.StatusMsg = &str
		response.MessageList = nil
		return
	}

	for _, msg := range msgList {
		apimsg := msg.ToApiMessage()
		response.MessageList = append(response.MessageList, apimsg)
		item := map[string]interface{}{
			"ID":        msg.ID,
			"FromID":    msg.FromID,
			"ToID":      msg.ToID,
			"Content":   msg.Content,
			"CreatedAt": msg.CreatedAt.UnixMilli(),
		}
		msg := RedisMsg{
			TYPE: ChatRecord,
			DATA: item,
		}
		ChanFromDB <- msg
	}
	response.StatusCode = 0
	str := "Get chat record successfully"
	response.StatusMsg = &str
}
