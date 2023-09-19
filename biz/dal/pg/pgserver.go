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
			log.Infoln("删除点赞缓存")
			go rdb.DelLikeVideo(msg.DATA)
		case Friend:
			log.Infoln("更新好友缓存-增加")
			go rdb.NewFriend(msg.DATA)
		case FriendDel:
			log.Infoln("更新好友缓存-删除")
			go rdb.DelFriend(msg.DATA)
		case FriendList:
			log.Infoln("更新好友列表缓存-增加")
			go rdb.UpdateFriendList(msg.DATA)
		case FriendListDel:
			log.Infoln("更新好友列表缓存-删除")
			go rdb.UpdateFriendListDel(msg.DATA)
		case ChatRecord:
			log.Infoln("更新聊天记录缓存")
			go rdb.UpdateChatRec(msg.DATA)
		}
	}
}
