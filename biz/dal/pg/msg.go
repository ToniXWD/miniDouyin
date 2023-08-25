package pg

import (
	"miniDouyin/biz/model/miniDouyin/api"
	"time"
)

const (
	DATE = "2006-01-02 15:04:05"
)

type DBMessage struct {
	ID      int64 `gorm:"primaryKey"`
	FromID  int64
	ToID    int64
	Content string
	Date    string
	Deleted int64 `gorm:"default:0"`
}

func (u *DBMessage) TableName() string {
	return "messages"
}

func (u *DBMessage) insert() bool {
	return DB.Create(u).Error == nil
}

func (u *DBMessage) ToApiMessage() (apimsg *api.Message) {
	apimsg = &api.Message{
		ID:         u.ID,
		ToUserID:   u.ToID,
		FromUserID: u.FromID,
		Content:    u.Content,
		CreateTime: u.Date,
	}
	return
}

func sendMsg(token string, toUerID int64, content string) bool {
	clientuser, _ := ValidateToken(token)
	t := time.Now()
	date := t.Format(DATE)
	msg := &DBMessage{
		FromID:  clientuser.ID,
		ToID:    toUerID,
		Content: content,
		Date:    date,
	}
	return msg.insert()
}
