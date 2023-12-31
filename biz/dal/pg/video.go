package pg

import (
	log "github.com/sirupsen/logrus"
	"miniDouyin/biz/dal/rdb"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
	"reflect"
	"strconv"
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

// 根据User的ID字段在数据库中查询
// 找到结果就填充整个结构体并返回T True
// 否则返回 False
func (v *DBVideo) QueryVideoByID() bool {
	result := DB.First(v, "ID = ?", v.ID)

	if result.Error != nil {
		return false
	}

	// 检查是否找到了记录
	return result.RowsAffected > 0
}

func (v *DBVideo) increaseComment(db *gorm.DB) bool {

	if v.ID <= 0 {
		// 视频iD不能小于等于0
		return false
	}
	res := db.Model(v).Where("ID = ?", v.ID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
	return res.Error == nil
}

// 被点赞数自增
// 需要保证ID有效
// 将当前结构体插入数据库，返回是否成功
// 需要提前保证该结构体有效
func (v *DBVideo) increaseFavorited(db *gorm.DB, num int64) *gorm.DB {
	return db.Model(v).Where("ID = ?", v.ID).Update("favorite_count", gorm.Expr("favorite_count + ?", num))
}

// 数据库模型转换为api的结构体
func (v *DBVideo) ToApiVideo(db *gorm.DB, clientUser *DBUser, islike bool) (*api.Video, error) {
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

	// 填充用户
	// 先尝试从Redis缓存获取作者
	dbuser, err := ID2DBUser(v.Author)
	if err != nil {
		return nil, err
	}

	av.Author, _ = dbuser.ToApiUser(clientUser)

	if clientUser == nil {
		// 如果客户端未登录，IsFavorite设置为false
		av.IsFavorite = false
		return av, nil
	}

	// islike设置为true表示调用者确认这个video已是被喜欢的状态
	if islike {
		av.IsFavorite = true
		return av, nil
	}

	// 判断用户是否对视频点过赞
	av.IsFavorite = IsVideoLikedByUser(clientUser.ID, v.ID)
	return av, nil
}

func (v *DBVideo) GetMinTimestamp() time.Time {
	var minViews time.Time
	DB.Model(v).Select("MIN(created_at)").Scan(&minViews)
	return minViews
}

func (v *DBVideo) InitSelfFromMap(uMap map[string]string) {
	reflectVal := reflect.ValueOf(v).Elem()

	for fieldName, fieldValue := range uMap {
		field := reflectVal.FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.Int, reflect.Int64:
				tmp, _ := strconv.ParseInt(fieldValue, 10, 64)
				field.SetInt(tmp)
			case reflect.String:
				field.SetString(fieldValue)
			case reflect.Struct:
				if field.Type() == reflect.TypeOf(time.Time{}) {
					// 如果字段类型是time.Time，尝试将字符串解析为时间
					t, err := utils.Str2Time(fieldValue)
					if err == nil {
						field.Set(reflect.ValueOf(t))
					}
				}
			}
		}
	}
}

// 更新视频缓存 和 用户发布视频的集合
// authorID 不为0表示新发布了视频
func (v *DBVideo) UpdateRedis(authorID int64) {
	// 更新视频缓存
	items1 := utils.StructToMap(v)
	msg1 := RedisMsg{
		TYPE: VideoInfo,
		DATA: items1,
	}
	ChanFromDB <- msg1

	if authorID == 0 {
		// authorID =0 表示视频被点赞或评论了，不需要更新用户发布视频的集合
		return
	}

	// 更新用户发布视频的集合
	msg2 := RedisMsg{
		TYPE: Publish,
		DATA: map[string]interface{}{
			"ID":     authorID,
			"Videos": []interface{}{v.ID},
		},
	}
	ChanFromDB <- msg2
}

// 返回至多30条视频列表
func GetNewVideoList(maxDate int64) (vlist []DBVideo, r_err error) {
	r_err = nil
	var dbv DBVideo
	videoNum := dbv.Count()

	if videoNum > 30 {
		videoNum = 30
	}
	//mintime := dbv.GetMinTimestamp()
	//
	//log.Debugf("mintime = %v\n", mintime)

	var cmp time.Time
	if maxDate <= 0 {
		cmp = time.Now()
	} else {
		cmp = time.UnixMilli(maxDate)
	}

	log.Debugf("cmp = %v\n", cmp)

	res := DB.Model(&DBVideo{}).Where("created_at <= ?", cmp).
		Order("ID desc").Limit(int(videoNum)).Find(&vlist)
	if res.Error != nil {
		r_err = res.Error
	}
	return
}

// 查询用户发布的视频列表，并更新到缓存集合
func GetUserVideoList(userID int64) (vlist []DBVideo, r_err error) {
	r_err = nil

	res := DB.Model(&DBVideo{}).Where("Author = ?", userID).
		Order("ID desc").Find(&vlist)
	if res.Error != nil {
		r_err = res.Error
		return
	}
	// 更新到缓存
	ids := make([]interface{}, len(vlist))
	for idx, item := range vlist {
		ids[idx] = item.ID
	}
	msg := RedisMsg{
		TYPE: Publish,
		DATA: map[string]interface{}{
			"ID":     userID,
			"Videos": ids,
		}}
	ChanFromDB <- msg
	return
}

func ID2VideoBy(ID int64) (*DBVideo, error) {
	dbv := &DBVideo{}
	// 先尝试缓存获取
	vMap, find := rdb.GetVideoById(strconv.FormatInt(ID, 10))
	if find {
		//	缓存命中
		dbv.InitSelfFromMap(vMap)
	} else {
		dbv.ID = ID
		find := dbv.QueryVideoByID()
		if !find {
			return nil, utils.ErrVideoNotExist
		}
		dbv.UpdateRedis(0)
	}
	return dbv, nil
}
