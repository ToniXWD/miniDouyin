/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-08-23 11:27:35
 * @LastEditTime: 2023-08-23 23:41:33
 * @version: 1.0
 */
package relation

import (
	"fmt"
	"miniDouyin/biz/model/miniDouyin/api"
)

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
