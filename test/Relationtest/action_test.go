package pg

import (
	"miniDouyin/biz/dal/pg"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// 测试能否正确处理关注请求
func TestAction_action(t *testing.T) {
	pg.Init()
	var action = &pg.DBAction{
		UserID:   4,
		FollowID: 5,
	}
	// 测试action是否能正确插入数据库
	assert.Equal(t, action.Insert(), nil)
	// 测试action是否能正确从数据库中删除
	pg.DB.Where("user_id = ? and follow_id = ?", 1, 2).Delete(&pg.DBAction{})
	user := &pg.DBUser{}
	pg.DB.Model(user).Where("ID = ?", 1).Update("follow_count", gorm.Expr("follow_count - 2"))
}
