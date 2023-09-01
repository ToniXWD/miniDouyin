package dal

import (
	"miniDouyin/biz/dal/pg"
	"miniDouyin/biz/dal/rdb"
)

func Init() {
	pg.Init()
	rdb.Init()
}
