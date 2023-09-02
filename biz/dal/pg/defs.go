/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-09-01 18:42:48
 * @LastEditTime: 2023-09-02 16:44:48
 * @version: 1.0
 */
package pg

var ChanFromDB chan RedisMsg

type RedisMsg struct {
	TYPE int
	DATA map[string]interface{}
}

const (
	UserInfo = iota // 用户信息
	UserFollowAdd
	UserFollowDel
	VideoInfo
	Publish // 用户上传视频的集合
	CommentCreate
	CommentDel
	Friend
	FriendList // 用户好友列表
	ChatRecord
)
