package pg

import (
	log "github.com/sirupsen/logrus"
	"miniDouyin/biz/dal/rdb"
)

func PGserver() {
	for {
		msg := <-ChanFromDB
		switch msg.TYPE {
		case UserInfo:
			log.Debugln("处理UserInfo消息...")
			go rdb.NewUser(msg.DATA)
		case UserFollowAdd:
			log.Debugln("处理UserFollowAdd消息...")
			go rdb.NewFollowRelation(msg.DATA)
		case UserFollowDel:
			log.Debugln("处理UserFollowAdd消息...")
			go rdb.DelFollowRelation(msg.DATA)
		}
	}
}
