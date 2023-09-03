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
	var dbuser DBUser

	// 先尝试从Redis缓存获取用户
	dbMap, find := rdb.GetUserById(v.Author)
	if find {
		// 如果缓存命中
		log.Debugln("ToApiVideo: 从缓存查询视频author记录成功")
		dbuser.InitSelfFromMap(dbMap)
	} else {
		// 否则需要重数据库加载用户
		res := DB.Model(&DBUser{}).First(&dbuser, "ID = ?", v.Author)
		if res.Error != nil {
			av.Author = nil
			return nil, utils.ErrVideoUserNotExist
		}
		// 发送消息更新用户缓存
		items := utils.StructToMap(&dbuser)
		msg := RedisMsg{
			TYPE: UserInfo,
			DATA: items,
		}
		ChanFromDB <- msg
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

	// 否则需要进行查询以判断前端用户是否喜欢该视频
	findres := &Like{}
	r := db.Model(&Like{}).First(findres, "user_id = ? AND video_id = ?", clientUser.ID, v.ID)
	if r.RowsAffected != 0 {
		// 如果找到了记录，设置IsFavorite为true
		av.IsFavorite = true
	}
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
