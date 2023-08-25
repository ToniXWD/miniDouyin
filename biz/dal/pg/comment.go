package pg

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"miniDouyin/biz/model/miniDouyin/api"
)

type Comment struct {
	gorm.Model
	VideoId int64
	UserId  int64
	Content string
}

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

func CreateComment(comment *Comment, ctx context.Context) (int64, error) {
	result := DB.WithContext(ctx).Create(comment)
	return int64(comment.ID), result.Error
}

func DeleteComment(commentId int64, ctx context.Context) error {
	var comment Comment
	session := DB.WithContext(ctx)
	// 先查询是否有此评论
	result := session.Where("id = ?", commentId).First(&comment)
	if result.Error != nil {
		return errors.New("del comment is not exist")
	}
	result = session.Where("id = ?", commentId).Delete(&Comment{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetCommentList(videoId int64, ctx context.Context) ([]*Comment, error) {
	var commentList []*Comment
	session := DB.WithContext(ctx)
	result := session.Model(Comment{}).Where(map[string]any{"video_id": videoId}).
		Order("created_at desc").
		Find(&commentList)
	if result.Error != nil {
		log.Println(result.Error)
		return commentList, errors.New("get comment list failed")
	}
	return commentList, nil
}

func GetCommentCnt(videoId int64, ctx context.Context) (int64, error) {
	var count int64
	session := DB.WithContext(ctx)
	result := session.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId}).
		Count(&count)
	if result.Error != nil {
		return 0, errors.New("find comments count failed")
	}
	return count, nil
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

type GetCommentListService struct {
	ctx context.Context
}

func NewGetCommentListService(ctx context.Context) *GetCommentListService {
	return &GetCommentListService{
		ctx: ctx,
	}
}

func (s *GetCommentListService) GetCommonList(req *api.CommentListRequest) ([]*api.Comment, error) {
	commentList := make([]*api.Comment, 0)
	comments, err := GetCommentList(req.VideoID, s.ctx)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return commentList, nil
	}

	clientUser, err := ValidateToken(req.Token)
	apiUser, _ := clientUser.ToApiUser(clientUser)
	for _, comm := range comments {
		commentList = append(commentList, &api.Comment{
			ID:         int64(comm.ID),
			User:       apiUser,
			Content:    comm.Content,
			CreateDate: comm.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return commentList, nil

}
