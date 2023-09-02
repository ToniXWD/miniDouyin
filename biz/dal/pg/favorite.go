package pg

import (
	"miniDouyin/utils"
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID        int64 `gorm:"primaryKey"`
	UserId    int64
	VideoId   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   gorm.DeletedAt `gorm:"default:NULL"`
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
	//  创建一条喜欢的记录后，还需要将被点赞者的获赞数+1
	find := curVideo.QueryVideoByID()
	if !find {
		tx.Rollback()
		return false
	}
	author := &DBUser{ID: curVideo.Author}
	find = author.QueryUserByID()
	if !find {
		tx.Rollback()
		return false
	}
	res = author.increaseFavorited(tx, 1)
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	//  创建一条喜欢的记录后，还需要将视频的总赞数+1
	res = curVideo.increaseFavorited(tx, 1)
	if res.Error != nil {
		tx.Rollback()
		return false
	}

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
	author := &DBUser{ID: curVideo.Author}
	find = author.QueryUserByID()
	if !find {
		tx.Rollback()
		return false
	}
	res = author.increaseFavorited(tx, -1)
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	//  创建一条喜欢的记录后，还需要将视频的总赞数-1
	res = curVideo.increaseFavorited(tx, -1)
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	// 提交事务
	res = tx.Commit()
	r_ans := res.Error == nil
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	return r_ans
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
	res := db.Model(&DBVideo{}).First(&dbv, "ID = ?", l.VideoId)
	return dbv, res.Error == nil
}
