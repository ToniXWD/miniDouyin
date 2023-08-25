package pg

import (
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"

	"gorm.io/gorm"
)

type DBUser struct {
	ID              int64  `gorm:"primaryKey"`
	Username        string `gorm:"unique"`
	Nickname        string
	Passwd          string
	FollowCount     int64 `gorm:"default:0"`
	FollowerCount   int64 `gorm:"default:0"`
	WorkCount       int64 `gorm:"default:0"`
	FavoriteCount   int64 `gorm:"default:0"`
	Token           string
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  int64          `gorm:"default:0"`
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

// 根据User的ID字段在数据库中查询
// 找到结果就填充整个结构体并返回T True
// 否则返回 False
func (u *DBUser) QueryUserByID() bool {
	result := DB.First(u, "ID = ?", u.ID)

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

// 将当前结构体插入数据库，返回是否成功
// 需要提前保证该结构体有效
func (u *DBUser) increaseWork(db *gorm.DB, num int64) bool {

	if u.Username == "" || u.Passwd == "" || u.ID == 0 {
		// 密码盒用户名不能为空
		return false
	}
	res := db.Model(u).Where("ID = ?", u.ID).Update("work_count", gorm.Expr("work_count + ?", num))
	return res.Error == nil
}

// 点赞数自增
// 将当前结构体插入数据库，返回是否成功
// 需要提前保证该结构体有效
func (u *DBUser) increaseFavorite(db *gorm.DB, num int64) *gorm.DB {
	return db.Model(u).Where("ID = ?", u.ID).Update("favorite_count", gorm.Expr("favorite_count + ?", num))
}

// 获赞数自增
// 将当前结构体插入数据库，返回是否成功
// 需要提前保证该结构体有效
func (u *DBUser) increaseFavorited(db *gorm.DB, num int64) *gorm.DB {
	return db.Model(u).Where("ID = ?", u.ID).Update("total_favorited", gorm.Expr("total_favorited + ?", num))
}

// 从数据库结构体转化为api的结构体
// IsFollow此时默认设置为false，后续需自行处理
// 表示查看当前用户的token
func (u *DBUser) ToApiUser(clientUser *DBUser) (apiuser *api.User, err error) {
	err = nil
	apiuser = &api.User{
		ID:              int64(u.ID),
		Name:            u.Username,
		FollowCount:     &u.FollowCount,
		FollowerCount:   &u.FollowerCount,
		IsFollow:        false,
		Avatar:          &u.Avatar,
		BackgroundImage: &u.BackgroundImage,
		Signature:       &u.Signature,
		TotalFavorited:  &u.TotalFavorited,
		WorkCount:       &u.WorkCount,
		FavoriteCount:   &u.FavoriteCount,
	}

	// 如果clientUser是空值，意味着前端没有登录，IsFollow默认false就可以返回
	if clientUser == nil {
		return
	}

	if clientUser.ID == apiuser.ID {
		// 自己一定是关注了自己的
		apiuser.IsFollow = true
		return
	}

	// 否则需要根据clientUser查询该用户是否被关注
	match := &DBAction{}
	e := DB.Model(&DBAction{}).Where("UserID = ? AND FollowID = ?", clientUser.ID, u.ID).First(match)
	if e.Error == nil {
		// 查询成功，match是有效的记录
		apiuser.IsFollow = true
		err = utils.ErrMathRealationFailed
	}
	return
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
	res := DB.First(&user, "ID = ?", request.UserID)

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

// 通过鉴权验证token
// 返回对应的user的指针和错误信息
func ValidateToken(token string) (*DBUser, error) {
	var user DBUser
	res := DB.First(&user, "Token = ?", token)
	if res.Error != nil {
		return nil, utils.ErrTokenVerifiedFailed
	}
	return &user, nil
}
