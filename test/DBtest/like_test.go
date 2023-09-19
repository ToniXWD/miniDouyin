package test

import (
	"fmt"
	"miniDouyin/biz/dal/pg"
	"miniDouyin/biz/dal/rdb"
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

func TestLike_VideoLikedByUser(t *testing.T) {
	pg.Init()
	rdb.Init()

	islike, err := rdb.IsVideoLikedById(8, 4)
	assert.Nil(t, err)
	assert.True(t, islike)

	islike, err = rdb.IsVideoLikedById(11, 4)
	assert.Nil(t, err)
	assert.True(t, islike)

	islike, err = rdb.IsVideoLikedById(12, 4)
	assert.Nil(t, err)
	assert.False(t, islike)
}
