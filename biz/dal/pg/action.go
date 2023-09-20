package pg

import (
	log "github.com/sirupsen/logrus"
	"miniDouyin/biz/dal/rdb"
	"miniDouyin/biz/model/miniDouyin/api"
	"strconv"

	"gorm.io/gorm"
)

type DBAction struct {
	ID       int64 `gorm:"primaryKey"`
	UserID   int64
	FollowID int64
	Deleted  gorm.DeletedAt `gorm:"default:NULL"`
}

func (u *DBAction) TableName() string {
	return "follows"
}

// 向数据库中删除当前结构体
func (u *DBAction) Delete() error {
	res := DB.Where("user_id = ? AND follow_id = ?", u.UserID, u.FollowID).Delete(u)
	return res.Error
}

// 向数据库中插入当前结构体
func (u *DBAction) Insert() error {
	res := DB.Create(u)
	return res.Error
}

// 根据请求类型进行关注或取消关注
func (u *DBAction) ifFollow(actiontype int64) error {
	var err error
	msg := RedisMsg{}
	if actiontype == 1 {
		// 关注
		err = u.Insert()
		if err == nil {
			// 粉丝用户数据更新
			fanUser, _ := ID2DBUser(u.UserID)
			// DB.Model(&DBUser{}).Where("ID = ?", u.UserID).Update("follow_count", gorm.Expr("follow_count + ?", 1))

			fanUser.FollowCount++
			fanUser.UpdateRedis()

			// 被关注者用户数据更新
			Followed, _ := ID2DBUser(u.FollowID)
			// DB.Model(&DBUser{}).Where("ID = ?", u.FollowID).Update("follower_count", gorm.Expr("follower_count + ?", 1))
			Followed.FollowerCount++
			Followed.UpdateRedis()

			// 判断对方是否关注了自己
			relation := &DBAction{}
			res := DB.Model(&DBAction{}).Find(relation, "user_id = ? AND follow_id = ?", u.FollowID, u.UserID)
			if res.RowsAffected != 0 {
				//查到对方也关注了自己
				err := AddFriend(u.FollowID, u.UserID)
				if err != nil {
					return err
				}
			}
		}
		msg.TYPE = UserFollowAdd

	} else if actiontype == 2 {
		// 取消关注
		err = u.Delete()
		if err == nil {
			user := &DBUser{}
			// 粉丝用户数据更新
			// DB.Model(user).Where("ID = ?", u.UserID).Update("follow_count", gorm.Expr("follow_count - ?", 1))
			// 被关注者用户数据更新
			// DB.Model(user).Where("ID = ?", u.FollowID).Update("follower_count", gorm.Expr("follower_count - ?", 1))

			// 判断对方是否关注了自己
			relation := &DBAction{}
			res := DB.Model(&DBAction{}).Find(relation, "user_id = ? AND follow_id = ?", u.FollowID, u.UserID)
			if res.RowsAffected != 0 {
				//查到对方也关注了自己
				err := DelFriend(u.FollowID, u.UserID)
				if err != nil {
					return err
				}
			}
		}
		msg.TYPE = UserFollowDel
	}

	// 发送消息更新缓存
	msg.DATA = map[string]interface{}{
		"UserID":   u.UserID,
		"FollowID": u.FollowID,
	}
	ChanFromDB <- msg
	return err
}

// 从关注请求返回新的DBAction结构体
func DBActionFromActionRequest(request *api.RelationActionRequest) *DBAction {
	clientUser, Terr := Token2DBUser(request.Token)
	if Terr != nil {
		return nil
	}
	return &DBAction{
		UserID:   int64(clientUser.ID),
		FollowID: request.ToUserID,
	}
}

// 从Redis获取关注列表
func GetFollowListFromRedis(response *api.RelationFollowListResponse, clientUser *DBUser) bool {
	str_ids, redisFind := rdb.GetFollowsIDList(clientUser.ID)
	if redisFind {
		// 缓存命中
		for _, id := range str_ids {
			ID, _ := strconv.Atoi(id)
			user, err := ID2DBUser(int64(ID))
			if err != nil {
				response.StatusCode = 1
				str := "Get follow list failed"
				response.StatusMsg = &str
				response.UserList = nil
				return false
			}
			apiUser, _ := user.ToApiUser(clientUser)
			response.UserList = append(response.UserList, apiUser)
		}
		log.Infoln("GetFollowListFromRedis: 从缓存完成关注列表获取")
		response.StatusCode = 0
		str := "Get follow list successfully"
		response.StatusMsg = &str
		return true
	}
	return false
}

// 从Redis获取粉丝列表
func GetFollowerListFromRedis(response *api.RelationFollowerListResponse, clientUser *DBUser) bool {
	str_ids, redisFind := rdb.GetFollowersIDList(clientUser.ID)
	if redisFind {
		// 缓存命中
		for _, id := range str_ids {
			ID, _ := strconv.Atoi(id)
			user, err := ID2DBUser(int64(ID))
			if err != nil {
				response.StatusCode = 1
				str := "Get follower list failed"
				response.StatusMsg = &str
				response.UserList = nil
				return false
			}

			apiUser, _ := user.ToApiUser(clientUser)
			response.UserList = append(response.UserList, apiUser)
		}
		log.Infoln("GetFollowerListFromRedis: 从缓存完成粉丝列表获取")
		response.StatusCode = 0
		str := "Get follow list successfully"
		response.StatusMsg = &str
		return true
	}
	return false
}
