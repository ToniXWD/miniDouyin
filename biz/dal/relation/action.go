package relation

import (
	"miniDouyin/biz/dal/pg"
	"miniDouyin/biz/model/miniDouyin/api"

	"gorm.io/gorm"
)

var DB *gorm.DB

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
	res := DB.Delete(u)
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
	if actiontype == 1 {
		// 关注
		err = u.Insert()
		if err == nil {
			user := &pg.DBUser{}
			DB.Model(user).Where("ID = ?", u.UserID).Update("follow_count", gorm.Expr("follow_count + ?", 1))
			DB.Model(user).Where("ID = ?", u.FollowID).Update("follower_count", gorm.Expr("follower_count + ?", 1))
			return err
		}

	} else if actiontype == 2 {
		// 取消关注
		err = u.Delete()
		if err == nil {
			user := &pg.DBUser{}
			DB.Model(user).Where("ID = ?", u.UserID).Update("follow_count", gorm.Expr("follow_count - ?", 1))
			DB.Model(user).Where("ID = ?", u.FollowID).Update("follower_count", gorm.Expr("follower_count - ?", 1))
			return err
		}
	}
	return err

}

// 从关注请求返回新的DBAction结构体
func DBActionFromActionRequest(request *api.RelationActionRequest) *DBAction {
	user, _ := pg.ValidateToken(request.Token)
	return &DBAction{
		UserID:   int64(user.ID),
		FollowID: request.ToUserID,
	}
}
