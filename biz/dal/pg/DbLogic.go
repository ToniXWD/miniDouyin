package pg

import (
	"fmt"
	"miniDouyin/biz/model/miniDouyin/api"
)

func DBUserLogin(request *api.UserLoginRequest, response *api.UserLoginResponse) {
	User := DBUserFromRequest(request)

	if User.QueryUser() {
		// user存在
		fmt.Printf("user = %+v\n", User)

		// 校验密码
		if User.Passwd != request.Password {
			response.StatusCode = 1
			response.StatusMsg = "Password wrong!"
			return
		}
		response.StatusCode = 0
		response.UserID = int64(User.ID)
		response.Token = User.Username + User.Passwd
		response.StatusMsg = "Login successfully!"
		return
	}
	response.StatusCode = 2
	response.StatusMsg = "User not exist!"
}
