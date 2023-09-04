package pg

import (
	log "github.com/sirupsen/logrus"
	"miniDouyin/biz/dal/rdb"
	"miniDouyin/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID        int64 `gorm:"primaryKey"`
	UserId    int64
	VideoId   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	//Deleted   gorm.DeletedAt `gorm:"default:NULL"` // 没有必要软删除
}

func (l *Like) TableName() string {
	return "favorited_videos"
}

// 将当前结构体插入数据库，返回是否成功
func (l *Like) insert(db *gorm.DB, clientUser *DBUser, curVideo *DBVideo) bool {
	if l.UserId == 0 || l.VideoId == 0 {
		// 用户和视频的id不能为空
		return false
	}

	tx := db.Begin()
	res := tx.Create(l)

	if res.Error != nil {
		tx.Rollback()
		return false
	}

	//  创建一条喜欢的记录后，还需要将对点赞者的总赞数+1
	res = clientUser.increaseFavorite(tx, 1)
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	// 更新点赞者缓存
	clientUser.FavoriteCount++
	clientUser.UpdateRedis()

	//  创建一条喜欢的记录后，还需要将被点赞者的获赞数+1
	find := curVideo.QueryVideoByID()
	if !find {
		tx.Rollback()
		return false
	}

	// 查询视频作者
	author, err := ID2DBUser(curVideo.Author)
	if err != nil {
		tx.Rollback()
		return false
	}

	res = author.increaseFavorited(tx, 1)
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	// 更新被点赞者缓存
	author.TotalFavorited++
	author.UpdateRedis()

	//  创建一条喜欢的记录后，还需要将视频的总赞数+1
	res = curVideo.increaseFavorited(tx, 1)
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	// 更新视频缓存
	curVideo.FavoriteCount++
	curVideo.UpdateRedis(0)

	// 提交事务
	res = tx.Commit()
	r_ans := res.Error == nil
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	return r_ans
}

// 将当前结构体删除，返回是否成功
func (l *Like) delete(db *gorm.DB, clientUser *DBUser, curVideo *DBVideo) bool {
	if l.UserId == 0 || l.VideoId == 0 {
		// 用户和视频的id不能为空
		return false
	}

	tx := db.Begin()
	res := db.Where("user_id = ? AND video_id = ?", l.UserId, l.VideoId).Delete(l)

	if res.Error != nil {
		tx.Rollback()
		return false
	}

	//  创建一条喜欢的记录后，还需要将对点赞者的总赞数-1
	res = clientUser.increaseFavorite(tx, -1)
	if res.Error != nil {
		tx.Rollback()
		return false
	}
	//  创建一条喜欢的记录后，还需要将被点赞者的获赞数-1
	find := curVideo.QueryVideoByID()
	if !find {
		tx.Rollback()
		return false
	}

	// 更新点赞者缓存
	clientUser.FavoriteCount--
	clientUser.UpdateRedis()

	// 查询视频作者
	author, err := ID2DBUser(curVideo.Author)
	if err != nil {
		tx.Rollback()
		return false
	}

	res = author.increaseFavorited(tx, -1)
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	// 更新被点赞者缓存
	author.TotalFavorited--
	author.UpdateRedis()

	//  创建一条喜欢的记录后，还需要将视频的总赞数-1
	res = curVideo.increaseFavorited(tx, -1)
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	// 更新视频缓存
	curVideo.FavoriteCount--
	curVideo.UpdateRedis(0)

	// 提交事务
	res = tx.Commit()
	rAns := res.Error == nil
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	return rAns
}

// 根据用户id查找喜欢的视频id
func (l *Like) QueryVideoByUser(db *gorm.DB) (dblist []DBVideo, find bool) {
	if l.UserId == 0 {
		return nil, false
	}
	var likelist []Like
	res := db.Model(&Like{}).Where("user_id = ?", l.UserId).
		Order("ID desc").Find(&likelist)

	// 检查是否找到了记录
	if res.RowsAffected > 0 {
		for _, item := range likelist {
			// 加入到缓存中
			items := utils.StructToMap(&item)
			msg := RedisMsg{
				TYPE: LikeCreate,
				DATA: items,
			}
			ChanFromDB <- msg

			dbv, _ := item.ToDBVideo(db)
			dblist = append(dblist, dbv)
		}
		return dblist, true
	}
	return nil, false
}

// 获取记录项中的视频
func (l *Like) ToDBVideo(db *gorm.DB) (dbv DBVideo, ans bool) {
	var find bool
	// 先尝试从缓存找到video
	vMap, find := rdb.GetVideoById(strconv.Itoa(int(l.VideoId)))
	if find {
		//	缓存命中
		dbv.InitSelfFromMap(vMap)
	} else {
		dbv.ID = l.VideoId
		find = dbv.QueryVideoByID()
	}
	return dbv, find
}

// 判断视频是否被用户点赞，先尝试查缓存，并更新缓存
func IsVideoLikedByUser(userID int64, videoID int64) bool {
	// 先从缓存判断
	islike, err := rdb.IsVideoLikedById(videoID, userID)
	if err == nil {
		//	缓存命中
		log.Infof("IsVideoLikedByUser：缓存命中：用户%v 与是否点赞了视频 %v? %t", userID, videoID, islike)
		return islike
	}
	// 否则查询数据库
	findres := &Like{}
	r := DB.Model(&Like{}).First(findres, "user_id = ? AND video_id = ?", userID, videoID)
	if r.RowsAffected != 0 {
		// 如果找到了记录，返回true
		// 更新缓存
		newRecord := &Like{
			UserId:  userID,
			VideoId: videoID,
		}
		items := utils.StructToMap(newRecord)
		msg := RedisMsg{
			TYPE: LikeCreate,
			DATA: items,
		}
		ChanFromDB <- msg
		log.Infoln("IsVideoLikedByUser：更新Like缓存")
		return true
	}
	return false
}
