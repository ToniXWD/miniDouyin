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
		case CommentCreate:
			log.Infoln("更新评论缓存")
			go rdb.NewComment(msg.DATA)
		case CommentDel:
			log.Infoln("删除评论缓存")
			go rdb.DelComment(msg.DATA)
		case LikeCreate:
			log.Infoln("更新点赞缓存")
			go rdb.NewLikeVideo(msg.DATA)
		case LikeDel:
			log.Infoln("删除评论缓存")
			go rdb.DelLikeVideo(msg.DATA)
		case FollowerCreate:
			log.Infoln("更新粉丝缓存")
			go rdb.NewFollower(msg.DATA)
		}
	}
}
