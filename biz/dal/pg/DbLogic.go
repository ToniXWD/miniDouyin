package pg

import (
	"fmt"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
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
	response.User = user.ToApiUser()
	response.StatusCode = 0
	str := "Get user information successfully!"
	response.StatusMsg = &str
}

// 处理视频流
func DBVideoFeed(request *api.FeedRequest, response *api.FeedResponse) {
	vlist, err := GetNewVideoList(*request.LatestTime)
	if err != nil {
		response.StatusCode = 1
		str := "Failed to get video list"
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
		newVideo, _ := video.ToApiVideo()
		response.VideoList = append(response.VideoList, newVideo)
	}
}

// 接受上传视频
func DBReceiveVideo(request *api.PublishActionRequest, response *api.PublishActionResponse) {

	//vName := utils.GetVideoName(request.Title)
}
