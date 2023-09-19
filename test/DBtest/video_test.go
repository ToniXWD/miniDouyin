package test

import (
	"miniDouyin/biz/dal/rdb"
	"testing"

	"miniDouyin/biz/dal/pg"

	"github.com/stretchr/testify/assert"
)

func TestVideo_Count(t *testing.T) {
	pg.Init()
	var v pg.DBVideo
	assert.Equal(t, v.Count(), int64(12))
}

func TestVideo_VMap2DBVideo(t *testing.T) {
	pg.Init()
	rdb.Init()
	var v pg.DBVideo
	vMap, find := rdb.GetVideoById("37")
	assert.Equal(t, true, find)
	v.InitSelfFromMap(vMap)
	t.Log(v)
}
