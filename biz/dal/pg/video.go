package pg

import (
	"fmt"
	"gorm.io/gorm"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
	"time"
)

type DBVideo struct {
	ID            int64 `gorm:"primaryKey"`
	Title         string
	Author        int64 // 外键关联到DBUser结构体的主键
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64 `gorm:"default:0"`
	CommentCount  int64 `gorm:"default:0"`
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

// 数据库模型转换为api的结构体
func (v *DBVideo) ToApiVideo() (*api.Video, error) {
	rPlayurl, _ := utils.Realurl(v.PlayUrl)
	rCoverurl, _ := utils.Realurl(v.CoverUrl)

	av := &api.Video{
		ID:            v.ID,
		PlayURL:       rPlayurl,
		CoverURL:      rCoverurl,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		IsFavorite:    true,
		Title:         v.Title,
	}

	var dbuser DBUser

	res := DB.Model(&DBUser{}).First(&dbuser, "ID = ?", v.Author)
	if res.Error != nil {
		av.Author = nil
		return nil, utils.ErrVideoUserNotExist
	}

	av.Author = dbuser.ToApiUser()

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
		maxDate = time.Now().Unix()
	}
	cmp := utils.I64ToTime(maxDate)
	res := DB.Model(&DBVideo{}).Where("created_at <= ?", cmp).
		Order("ID desc").Limit(int(videoNum)).Find(&vlist)
	if res.Error != nil {
		r_err = res.Error
	}
	return
}
