package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"miniDouyin/biz/dal/pg"
)

func TestVideo_Count(t *testing.T) {
	pg.Init()
	var v pg.DBVideo
	assert.Equal(t, v.Count(), int64(1))
}
