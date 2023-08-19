package pg

import (
	"miniDouyin/biz/model/miniDouyin/api"

	"gorm.io/gorm"
)

type DBUser struct {
	ID              int    `gorm:"primaryKey"`
	Username        string `gorm:"unique"`
	Nickname        string
	Passwd          string
	FollowCount     int `gorm:"default:0"`
	FollowerCount   int `gorm:"default:0"`
	WorkCount       int `gorm:"default:0"`
	FavoriteCount   int `gorm:"default:0"`
	Token           string
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  int            `gorm:"default:0"`
	Deleted         gorm.DeletedAt `gorm:"default:NULL"`
}

func (u *DBUser) TableName() string {
	return "users"
}

func (u *DBUser) QueryUser() bool {
	if u.Username == "" || u.Passwd == "" {
		return false
	}
	result := DB.First(u, "Username = ?", u.Username)

	if result.Error != nil {
		return false
	}

	// 检查是否找到了记录
	return result.RowsAffected > 0
}

// 从注册请求构造新用户
func DBUserFromRequest(request *api.UserLoginRequest) DBUser {
	return DBUser{
		Username: request.Username,
		Passwd:   request.Password,
		Token:    request.Username + request.Password,
	}
}
