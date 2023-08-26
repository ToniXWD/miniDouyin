package pg

import (
	"fmt"
	"miniDouyin/biz/model/miniDouyin/api"
	"time"

	"gorm.io/gorm"
)

const (
	DATE = "2006-01-02 15:04:05"
)

type DBMessage struct {
	ID        int64 `gorm:"primaryKey"`
	FromID    int64
	ToID      int64
	Content   string
	CreatedAt time.Time
	Deleted   gorm.DeletedAt `gorm:"default:NULL"`
}

func (u *DBMessage) TableName() string {
	return "messages"
}

func (u *DBMessage) insert() bool {
	res := DB.Create(u)
	return res.Error == nil
}

func (u *DBMessage) ToApiMessage() (apimsg *api.Message) {
	time := u.CreatedAt.UnixMilli()
	fmt.Println("传出的时间戳为:", time)
	apimsg = &api.Message{
		ID:         u.ID,
		ToUserID:   u.ToID,
		FromUserID: u.FromID,
		Content:    u.Content,
		CreateTime: &time,
	}
	return
}

func sendMsg(token string, toUerID int64, content string) bool {
	clientuser, _ := ValidateToken(token)
	msg := &DBMessage{
		FromID:    clientuser.ID,
		ToID:      toUerID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	return msg.insert()
}
