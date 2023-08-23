package relation

import (
	"miniDouyin/biz/dal/relation"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试能否正确处理关注请求
func TestAction_action(t *testing.T) {
	relation.Init()
	var action = &relation.DBAction{
		UserID:   4,
		FollowID: 5,
	}
	// 测试action是否能正确插入数据库
	assert.Equal(t, action.Insert(), nil)
	// 测试action是否能正确从数据库中删除
	// relation.DB.Where("user_id = ? and follow_id = ?", 1, 2).Delete(&relation.DBAction{})
	// user := &pg.DBUser{}
	// relation.DB.Model(user).Where("ID = ?", 1).Update("follow_count", gorm.Expr("follow_count - 2"))
}
