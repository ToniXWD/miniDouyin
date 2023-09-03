package pg

import (
	"miniDouyin/biz/dal/rdb"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        int64 `gorm:"primaryKey"`
	VideoId   int64
	UserId    int64
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   gorm.DeletedAt `gorm:"default:NULL"`
}

func (v *Comment) TableName() string {
	return "comments"
}

// 数据库模型转换为api的结构体
func (v *Comment) ToApiComment(cUser *DBUser, clientUser *DBUser) (*api.Comment, error) {
	ac := &api.Comment{
		ID:         v.ID,
		User:       nil,
		Content:    v.Content,
		CreateDate: v.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 填充评论用户信息
	capiUser, err := cUser.ToApiUser(clientUser)
	if err != nil {
		return ac, err
	}
	ac.User = capiUser

	return ac, nil
}

// 发布评论
func (v *Comment) CreateComment() (int64, error) {
	result := DB.Create(v)
	return v.ID, result.Error
}

// 删除评论
func DeleteComment(commentId int64) (*Comment, error) {
	comment := Comment{}
	// 先查询是否有此评论
	result := DB.Where("id = ?", commentId).First(&comment)
	if result.Error != nil {
		return nil, utils.ErrDelCommentNotExist
	}

	// 如果评论存在则删除此条评论
	result = DB.Where("id = ?", commentId).Delete(&Comment{})
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

// 根据 video_id获取评论列表
func GetDBCommentList(v_id int64) (clist []Comment, err error) {
	err = nil
	res := DB.Model(&Comment{}).Where("video_id = ?", v_id).
		Order("ID").Find(&clist)
	if res.Error != nil {
		err = utils.ErrGetCommentListFailed
	}
	return
}

func NewGetCommentListService(v_id int64, token string) (clist []*api.Comment, r_err error) {
	r_err = nil
	var clientUser *DBUser = nil
	if token != "" {
		// token不为空，表示客户端已经登录
		// 校验token
		// 尝试从缓存查询用户
		uMap, find := rdb.GetUserByToken(token)
		if find {
			// 缓存命中
			clientUser.InitSelfFromMap(uMap)
		} else {
			// 缓存未命中，从数据库查
			clientUser, r_err = ValidateToken(token)
			if r_err != nil {
				return nil, r_err
			}
			// 更新缓存
			items := utils.StructToMap(clientUser)
			msg := RedisMsg{
				TYPE: UserInfo,
				DATA: items,
			}
			ChanFromDB <- msg
		}
	}

	// 校验视频id合法性
	v := &DBVideo{}
	res := DB.Model(v).First(v, "ID = ?", v_id)
	if res.Error != nil {
		return nil, utils.ErrVideoNotExist
	}

	// 如果视频id有效再获取评论列表
	// 缓存未命中，从数据库查
	// TODO: 先尝试从缓存查找视频评论列表
	// TODO：如果从缓存找到了视频评论列表，再从获取到的视频id查询评论（也是先尝试缓存查，再数据库查）
	cDBlist, err := GetDBCommentList(v_id)
	if err != nil {
		return nil, err
	}
	// 将评论列表格式进行转换
	for _, dbcomment := range cDBlist {
		// 更新缓存
		items := utils.StructToMap(&dbcomment)
		msg := RedisMsg{
			TYPE: CommentCreate,
			DATA: items,
		}
		ChanFromDB <- msg

		cUser := &DBUser{ID: dbcomment.UserId}
		if !cUser.QueryUserByID() {
			return nil, utils.ErrUserNotFound
		}
		ac, err := dbcomment.ToApiComment(cUser, clientUser)
		if err != nil {
			return nil, utils.ErrGetCommentListFailed
		}
		clist = append(clist, ac)
	}
	return
}

// 根据评论请求封装发布评论
func dbCreateComment(req *api.CommentActionRequest, userId int64) (*Comment, error) {
	comm := &Comment{
		VideoId: req.VideoID,
		UserId:  userId,
		Content: *req.CommentText,
	}
	_, err := comm.CreateComment()
	if err != nil {
		return nil, err
	}
	return comm, nil
}

/*
	评论模块架构设计
		发送评论
		当接收用户
*/
