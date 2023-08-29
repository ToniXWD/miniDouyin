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
			log.Infoln("更新user缓存...")
			go rdb.NewUser(msg.DATA)
		case UserFollowAdd:
			log.Infoln("更新关注缓存...")
			go rdb.NewFollowRelation(msg.DATA)
		case UserFollowDel:
			log.Infoln("更新取关缓存...")
			go rdb.DelFollowRelation(msg.DATA)
		case VideoInfo:
			log.Infoln("更新视频缓存...")
			go rdb.NewVideo(msg.DATA)
		case Publish:
			log.Infoln("更新用户发布视频列表缓存...")
			go rdb.NewPublish(msg.DATA)
		}
	}
}
