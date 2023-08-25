/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-08-25 16:15:08
 * @LastEditTime: 2023-08-25 16:58:26
 * @version: 1.0
 */
package pg

import "gorm.io/gorm"

type DBfriend struct {
	ID       int64 `gorm:"primaryKey"`
	UserID   int64
	FriendID int64
	Deleted  gorm.DeletedAt `gorm:"default:NULL"`
}

func (u *DBfriend) TableName() string {
	return "friends"
}
