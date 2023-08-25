package pg

import (
	"context"
	"errors"
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

// type DBVideo struct {
// 	ID            int64 `gorm:"primaryKey"`
// 	Title         string
// 	Author        int64 // 外键关联到DBUser结构体的主键
// 	PlayUrl       string
// 	CoverUrl      string `gorm:"default:'defaults/douyin.jpg'"`
// 	FavoriteCount int64  `gorm:"default:0"`
// 	CommentCount  int64  `gorm:"default:0"`
// 	CreatedAt     time.Time
// 	Deleted       gorm.DeletedAt `gorm:"default:NULL"`
// }

func (Comment) TableName() string {
	return "comments"
}

// api.Comment
// type Comment struct {
//	// 视频评论id
//	ID int64 `thrift:"id,1,required" form:"id,required" json:"id,required" query:"id,required"`
//	// 评论用户信息
//	User *User `thrift:"user,2,required" form:"user,required" json:"user,required" query:"user,required"`
//	// 评论内容
//	Content string `thrift:"content,3,required" form:"content,required" json:"content,required" query:"content,required"`
//	// 评论发布日期，格式 mm-dd
//	CreateDate string `thrift:"create_date,4,required" form:"create_date,required" json:"create_date,required" query:"create_date,required"`
// }

// 数据库模型转换为api的结构体
func (v *Comment) ToApiComment(cUser *DBUser, clientUser *DBUser) (*api.Comment, error) {
	ac := &api.Comment{
		ID:         v.ID,
		User:       nil,
		Content:    v.Content,
		CreateDate: v.CreatedAt.Format("mm-dd"),
	}

	// 填充评论用户信息
	capiUser, err := cUser.ToApiUser(clientUser)
	if err != nil {
		return ac, err
	}
	ac.User = capiUser

	return ac, nil
}

func CreateComment(comment *Comment, ctx context.Context) (int64, error) {
	result := DB.Create(comment)
	return int64(comment.ID), result.Error
}

func DeleteComment(commentId int64, ctx context.Context) error {
	comment := Comment{}
	// 先查询是否有此评论
	result := DB.Where("id = ?", commentId).First(&comment)
	if result.Error != nil {
		return errors.New("del comment is not exist")
	}
	result = DB.Where("id = ?", commentId).Delete(&Comment{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type ActionCommentService struct {
	ctx context.Context
}

func NewActionCommentService(ctx context.Context) *ActionCommentService {
	return &ActionCommentService{
		ctx: ctx,
	}
}
func (s *ActionCommentService) CreateComment(req *api.CommentActionRequest, userId int64) (*Comment, error) {
	comm := &Comment{
		VideoId: req.VideoID,
		UserId:  userId,
		Content: *req.CommentText,
	}
	_, err := CreateComment(comm, s.ctx)
	if err != nil {
		return nil, err
	}
	return comm, nil
}

func (s *ActionCommentService) DeleteComment(req *api.CommentActionRequest) error {
	err := DeleteComment(*req.CommentID, s.ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewGetCommentListService(v_id int64, token string) (clist []*api.Comment, r_err error) {
	r_err = nil
	var clientUser *DBUser = nil
	if token != "" {
		// token不为空，表示客户端已经登录
		// 校验token
		clientUser, r_err = ValidateToken(token)
		if r_err != nil {
			return nil, r_err
		}
	}

	// 校验视频id合法性
	v := &DBVideo{}
	res := DB.Model(v).First(v, "ID = ?", v_id)
	if res.Error != nil {
		return nil, utils.ErrVideoNotExist
	}

	// 如果视频id有效再获取评论列表
	cDBlist, err := GetDBCommentList(v_id)
	if err != nil {
		return nil, err
	}
	// 将视频列表格式进行转换
	for _, dbcomment := range cDBlist {
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

func GetDBCommentList(v_id int64) (clist []Comment, err error) {
	err = nil
	res := DB.Model(&Comment{}).Where("video_id = ?", v_id).
		Order("ID").Find(&clist)
	if res.Error != nil {
		err = utils.ErrGetCommentListFailed
	}
	return
}
