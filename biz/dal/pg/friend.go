/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-08-25 16:15:08
 * @LastEditTime: 2023-08-25 16:58:26
 * @version: 1.0
 */
package pg

import (
	"miniDouyin/biz/model/miniDouyin/api"
)

type DBfriend struct {
	ID       int64 `gorm:"primaryKey"`
	UserID   int64
	FriendID int64
}

func (u *DBfriend) TableName() string {
	return "friends"
}

func UpdateRedis(id1 int64, id2 int64, Type1 int, Type2 int) {
	// 更新缓存一个好友关系
	item1 := map[string]interface{}{
		"FromID":  id1,
		"ToID":    id2,
		"Message": "",
		"MsgType": 0,
	}
	msg1 := RedisMsg{
		TYPE: Type1,
		DATA: item1,
	}
	ChanFromDB <- msg1

	// 更新另一个好友关系
	item2 := map[string]interface{}{
		"FromID":  id2,
		"ToID":    id1,
		"Message": "",
		"MsgType": 0,
	}
	msg2 := RedisMsg{
		TYPE: Type1,
		DATA: item2,
	}
	ChanFromDB <- msg2

	// 更新集合关系1
	item3 := map[string]interface{}{
		"ID":     id1,
		"Friend": id2,
	}
	msg3 := RedisMsg{
		TYPE: Type2,
		DATA: item3,
	}
	ChanFromDB <- msg3

	// 更新集合关系2
	item4 := map[string]interface{}{
		"ID":     id2,
		"Friend": id1,
	}
	msg4 := RedisMsg{
		TYPE: Type2,
		DATA: item4,
	}
	ChanFromDB <- msg4

}

// 插入2条朋友关系
func AddFriend(id1 int64, id2 int64) error {
	f1 := &DBfriend{}
	f1.UserID = id1
	f1.FriendID = id2
	res := DB.Model(f1).Create(f1)
	if res.Error != nil {
		return res.Error
	}

	f2 := &DBfriend{}
	f2.FriendID = id1
	f2.UserID = id2
	res = DB.Model(f2).Create(f2)
	if res.Error != nil {
		return res.Error
	}

	UpdateRedis(id1, id2, Friend, FriendList)

	return nil
}

// 删除2条朋友关系
func DelFriend(id1 int64, id2 int64) error {
	dbf := &DBfriend{}
	res := DB.Model(&DBfriend{}).Where("user_id = ? AND friend_id = ?", id1, id2).Delete(dbf)
	if res.Error != nil {
		return res.Error
	}
	res = DB.Model(&DBfriend{}).Where("user_id = ? AND friend_id = ?", id2, id1).Delete(dbf)
	if res.Error != nil {
		return res.Error
	}

	UpdateRedis(id1, id2, FriendDel, FriendListDel)

	return nil
}

// 将user转化为friend
func apiUser2apiFriend(user *api.User, client *DBUser) (friend *api.FriendUser, msg DBMessage) {
	friend = new(api.FriendUser)
	// 先进行公共部分赋值转化
	friend.IsFollow = user.IsFollow
	friend.WorkCount = user.WorkCount
	friend.ID = user.ID
	friend.FavoriteCount = user.FavoriteCount
	friend.BackgroundImage = user.BackgroundImage
	friend.Signature = user.Signature
	friend.TotalFavorited = user.TotalFavorited
	friend.FollowerCount = user.FollowerCount
	friend.FollowCount = user.FollowCount
	friend.Avatar = user.Avatar
	friend.Name = user.Name

	// 查询二人最新的消息记录

	err := DB.Where("((from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?))", client.ID, user.ID, user.ID, client.ID).Order("ID desc").Limit(1).First(&msg)
	if err.Error != nil {
		return
	}

	message := msg.Content
	friend.Message = &message

	if msg.ToID == client.ID {
		friend.MsgType = 0
	} else {
		friend.MsgType = 1
	}
	return
}

//// 从Redis获取粉丝列表
//func GetFriendListFromRedis(response *api.RelationFriendListResponse, clientUser *DBUser) bool {
//	str_ids, redisFind := rdb.GetFriendList(clientUser.ID)
//	if redisFind {
//		// 缓存命中
//		for _, id := range str_ids {
//			ID, _ := strconv.Atoi(id)
//			user, err := ID2DBUser(int64(ID))
//			if err != nil {
//				response.StatusCode = 1
//				str := "Get follower list failed"
//				response.StatusMsg = &str
//				response.UserList = nil
//				return false
//			}
//
//			apiUser, _ := user.ToApiUser(clientUser)
//			response.UserList = append(response.UserList, apiUser)
//		}
//		log.Infoln("GetFollowerListFromRedis: 从缓存完成粉丝列表获取")
//		response.StatusCode = 0
//		str := "Get follow list successfully"
//		response.StatusMsg = &str
//		return true
//	}
//	return false
//}
