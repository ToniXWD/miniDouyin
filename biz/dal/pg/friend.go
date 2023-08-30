/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-08-25 16:15:08
 * @LastEditTime: 2023-08-25 16:58:26
 * @version: 1.0
 */
package pg

import (
	"gorm.io/gorm"
	"miniDouyin/biz/model/miniDouyin/api"
)

type DBfriend struct {
	ID       int64 `gorm:"primaryKey"`
	UserID   int64
	FriendID int64
	Deleted  gorm.DeletedAt `gorm:"default:NULL"`
}

func (u *DBfriend) TableName() string {
	return "friends"
}

// 将user转化为friend
func apiUser2apiFriend(user *api.User, client *DBUser) (friend *api.FriendUser) {
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
	var msg DBMessage
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
