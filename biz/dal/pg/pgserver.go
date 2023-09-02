/*
 * @Description:a
 * @Author: Zjy
 * @Date: 2023-09-01 18:42:48
 * @LastEditTime: 2023-09-02 18:58:37
 * @version: 1.0
 */
package pg

import (
	"miniDouyin/biz/dal/rdb"

	log "github.com/sirupsen/logrus"
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
		case Friend:
			log.Infoln("更新好友缓存")
			go rdb.NewFriend(msg.DATA)
		case FriendList:
			log.Infoln("更新好友列表缓存")
			go rdb.UpdateFriendList(msg.DATA)
		case ChatRecord:
			log.Infoln("更新聊天记录缓存")
			go rdb.UpdateChatRec(msg.DATA)
		}
	}
}
