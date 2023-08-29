package pg

import (
	"fmt"
	"miniDouyin/biz/dal/rdb"
)

func PGserver() {
	for {
		msg := <-ChanFromDB
		switch msg.TYPE {
		case UserInfo:
			fmt.Println("处理UserInfo消息...")
			go rdb.NewUser(msg.DATA)
		case UserFollowAdd:
			fmt.Println("处理UserFollowAdd消息...")
			go rdb.NewFollowRelation(msg.DATA)
		case UserFollowDel:
			fmt.Println("处理UserFollowAdd消息...")
			go rdb.DelFollowRelation(msg.DATA)
		}
	}
}
