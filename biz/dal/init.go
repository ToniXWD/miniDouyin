package dal

import (
	"miniDouyin/biz/dal/pg"
	"miniDouyin/biz/dal/relation"
)

func Init() {
	pg.Init()
	relation.Init()
}
