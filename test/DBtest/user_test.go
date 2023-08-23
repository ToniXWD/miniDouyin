package test

import (
	"fmt"
	"miniDouyin/biz/dal/pg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Token(t *testing.T) {
	pg.Init()

	token := "toni123456"
	user, _ := pg.ValidateToken(token)
	assert.NotNil(t, user)
	assert.Equal(t, user.Username, "toni")
	assert.Equal(t, user.Passwd, "123456")
}

func TestUser_Work(t *testing.T) {
	pg.Init()
	var users []pg.DBUser
	pg.DB.Where("ID > ?", 0).Find(&users)
	for idx, user := range users {
		fmt.Printf("第%v个结果:\n", idx)
		fmt.Printf("\t%v\n\n", user)
	}
}
