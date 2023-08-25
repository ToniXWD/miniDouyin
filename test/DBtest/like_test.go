package test

import (
	"fmt"
	"miniDouyin/biz/dal/pg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLike_List(t *testing.T) {
	pg.Init()
	user_id := 4
	like := &pg.Like{
		UserId: int64(user_id),
	}

	list, find := like.QueryVideoByUser(pg.DB)
	assert.Equal(t, true, find)
	fmt.Printf("%+v", list)
}
