package dal

import (
	"miniDouyin/biz/dal/pg"
	"miniDouyin/biz/dal/rdb"
)

func Init() {
	rdb.Init()
	pg.Init()
}
