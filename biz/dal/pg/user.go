package pg

import (
	"gorm.io/gorm"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
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

// 根据User的Username字段在数据库中查询
// 找到结果就填充整个结构体并返回T True
// 否则返回 False
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

// 将当前结构体插入数据库，返回是否成功
func (u *DBUser) insert() bool {
	if u.Username == "" || u.Passwd == "" {
		// 密码盒用户名不能为空
		return false
	}

	// 设置token
	if u.Token == "" {
		u.Token = u.Username + u.Passwd
	}

	res := DB.Create(u)

	return res.Error == nil
}

// 从数据库结构体转化为api的结构体
// IsFollow此时默认设置为true，后续需自行处理
func (u *DBUser) ToApiUser() *api.User {
	return &api.User{
		ID:              int64(u.ID),
		Name:            u.Username,
		FollowCount:     int64(u.FollowCount),
		FollowerCount:   int64(u.FollowerCount),
		IsFollow:        true,
		Avatar:          u.Avatar,
		BackgroundImage: u.BackgroundImage,
		Signature:       u.Signature,
		TotalFavorited:  int64(u.TotalFavorited),
		WorkCount:       int64(u.WorkCount),
		FavoriteCount:   int64(u.FavoriteCount),
	}
}

// 从登录请求构造新用户
func DBUserFromLoginRequest(request *api.UserLoginRequest) DBUser {
	return DBUser{
		Username: request.Username,
		Passwd:   request.Password,
		Token:    request.Username + request.Password,
	}
}

// 从注册请求构造新用户
func DBUserFromRegisterRequest(request *api.UserRegisterRequest) DBUser {
	return DBUser{
		Username: request.Username,
		Passwd:   request.Password,
		Token:    request.Username + request.Password,
	}
}

// 从获取用户信息请求请求构造新用户
func DBGetUser(request *api.UserRequest) (*DBUser, error) {
	var user DBUser
	res := DB.First(&user, request.UserID)

	if res.Error != nil {
		// 没有找到记录
		return nil, utils.ErrUserNotFound
	}

	// 找到后，与token进行比对
	if user.Token != request.Token {
		return nil, utils.ErrWrongToken
	}

	return &user, nil
}
