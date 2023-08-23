package test

import (
	"testing"

	"miniDouyin/biz/dal/pg"

	"github.com/stretchr/testify/assert"
)

func TestVideo_Count(t *testing.T) {
	pg.Init()
	var v pg.DBVideo
	assert.Equal(t, v.Count(), int64(12))
}
