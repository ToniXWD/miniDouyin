package pg

import (
	"fmt"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
	"time"

	"gorm.io/gorm"
)

type DBVideo struct {
	ID            int64 `gorm:"primaryKey"`
	Title         string
	Author        int64 // 外键关联到DBUser结构体的主键
	PlayUrl       string
	CoverUrl      string `gorm:"default:'defaults/douyin.jpg'"`
	FavoriteCount int64  `gorm:"default:0"`
	CommentCount  int64  `gorm:"default:0"`
	CreatedAt     time.Time
	Deleted       gorm.DeletedAt `gorm:"default:NULL"`
}

func (v *DBVideo) TableName() string {
	return "videos"
}

// 统计videos数量
func (v *DBVideo) Count() int64 {
	var videoNum int64
	DB.Model(v).Count(&videoNum)
	return videoNum
}

// 将当前结构体插入数据库，返回是否成功
func (v *DBVideo) insert(db *gorm.DB) bool {
	if v.PlayUrl == "" || v.Author == 0 {
		// PlayUrl 和 PlayUrl不能为空
		return false
	}

	res := db.Create(v)

	return res.Error == nil
}

func (v *DBVideo) increaseComment(db *gorm.DB) bool {

	if v.ID <= 0 {
		// 视频iD不能小于等于0
		return false
	}
	res := db.Model(v).Where("ID = ?", v.ID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
	return res.Error == nil
}

// 数据库模型转换为api的结构体
func (v *DBVideo) ToApiVideo(clientUser *DBUser) (*api.Video, error) {
	rPlayurl := utils.Realurl(v.PlayUrl)
	rCoverurl := utils.Realurl(v.CoverUrl)

	av := &api.Video{
		ID:            v.ID,
		PlayURL:       rPlayurl,
		CoverURL:      rCoverurl,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		IsFavorite:    false,
		Title:         v.Title,
	}

	var dbuser DBUser

	res := DB.Model(&DBUser{}).First(&dbuser, "ID = ?", v.Author)
	if res.Error != nil {
		av.Author = nil
		return nil, utils.ErrVideoUserNotExist
	}

	av.Author, _ = dbuser.ToApiUser(clientUser)

	return av, nil
}

func (v *DBVideo) GetMinTimestamp() time.Time {
	var minViews time.Time
	DB.Model(v).Select("MIN(created_at)").Scan(&minViews)
	return minViews
}

// 返回至多30条视频列表
func GetNewVideoList(maxDate int64) (vlist []DBVideo, r_err error) {
	r_err = nil
	var dbv DBVideo
	videoNum := dbv.Count()

	if videoNum > 30 {
		videoNum = 30
	}
	mintime := dbv.GetMinTimestamp()

	fmt.Printf("mintime = %v\n", mintime)

	if maxDate <= 0 {
		maxDate = time.Now().Unix() * 1000
	}
	cmp := utils.I64ToTime(maxDate)

	fmt.Printf("cmp = %v\n", cmp)

	res := DB.Model(&DBVideo{}).Where("created_at <= ?", cmp).
		Order("ID desc").Limit(int(videoNum)).Find(&vlist)
	if res.Error != nil {
		r_err = res.Error
	}
	return
}

func GetUserVideoList(userID int64) (vlist []DBVideo, r_err error) {
	r_err = nil

	res := DB.Model(&DBVideo{}).Where("Author = ?", userID).
		Order("ID desc").Find(&vlist)
	if res.Error != nil {
		r_err = res.Error
	}
	return
}
