package pg

import (
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
