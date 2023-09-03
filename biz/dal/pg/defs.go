package pg

var ChanFromDB chan RedisMsg

type RedisMsg struct {
	TYPE int
	DATA map[string]interface{}
}

const (
	UserInfo = iota // 用户信息
	UserFollowAdd
	UserFollowDel
	VideoInfo
	Publish // 用户上传视频的集合
	CommentCreate
	CommentDel
	LikeCreate
	LikeDel
)
