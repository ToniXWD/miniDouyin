package pg

import (
	"context"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"

	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserId  int64
	VideoId int64
}

func (f *Like) TableName() string {
	return "favorited_videos"
}

// AddFavorite 点赞操作
func AddFavorite(like *Like, ctx context.Context) (int64, error) {
	result := DB.WithContext(ctx).Create(like)
	return int64(like.ID), result.Error
}

// CancelFavorite 取消点赞操作
func CancelFavorite(like *Like, ctx context.Context) error {
	result := DB.WithContext(ctx).Where("user_id = ? AND video_id = ?", like.UserId, like.VideoId).Delete(like)
	return result.Error
}

//

// // ListFavorite 返回userid用户点赞的视频列表
func ListFavorite(userId int64, ctx context.Context) ([]*DBVideo, []int64, error) {
	data := make([]*Like, 0)
	if err := DB.WithContext(ctx).Where("user_id = ?", userId).Find(&data).Error; err != nil {
		return nil, nil, err
	}
	if len(data) == 0 {
		return make([]*DBVideo, 0), nil, nil
	}

	videoIdList := make([]int64, len(data))
	for i, like := range data {
		videoIdList[i] = like.VideoId
	}

	countList, err := MCountFavorite(videoIdList, ctx)
	if err != nil {
		return nil, nil, err
	}

	DBVideos, err := MGetVideos(videoIdList, ctx)
	if err != nil {
		return nil, countList, err
	}

	return DBVideos, countList, err

}

// MCountFavorite 统计videoId点赞数量
func MCountFavorite(videoIdList []int64, ctx context.Context) ([]int64, error) {
	countList := make([]int64, len(videoIdList))
	for i, videoId := range videoIdList {
		if err := DB.WithContext(ctx).Model(&Like{}).Where("video_id = ?", videoId).Count(&countList[i]).Error; err != nil {
			return nil, err
		}
	} //
	return countList, nil
}

// MGetVideos 通过 videoID 获取对应 DBVideo
func MGetVideos(videoIdList []int64, ctx context.Context) ([]*DBVideo, error) {
	data := make([]*DBVideo, 0)
	if err := DB.WithContext(ctx).Where("id in ?", videoIdList).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

type FavoriteActionService struct {
	ctx context.Context
}

func NewFavoriteActionService(ctx context.Context) *FavoriteActionService {
	return &FavoriteActionService{ctx: ctx}
}

func (s *FavoriteActionService) AddFavorite(req *api.FavoriteActionRequest, userId int64) error {
	_, err := AddFavorite(&Like{
		UserId:  userId,
		VideoId: req.VideoID,
	}, s.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *FavoriteActionService) CancelFavorite(req *api.FavoriteActionRequest, userId int64) error {
	err := CancelFavorite(&Like{
		UserId:  userId,
		VideoId: req.VideoID,
	}, s.ctx)
	if err != nil {
		return err
	}
	return nil
}

type FavoriteListService struct {
	ctx context.Context
}

// NewFavoriteListService creates a new FavoriteListService
func NewFavoriteListService(ctx context.Context) *FavoriteListService {
	return &FavoriteListService{
		ctx: ctx,
	}
}

func (s *FavoriteListService) ListFavorite(req *api.FavoriteListRequest, userId int64) ([]*api.Video, error) {
	clientUser, _ := ValidateToken(req.Token)
	videos, countList, err := ListFavorite(req.UserID, s.ctx)
	if err != nil {
		return nil, utils.ErrWrongParam
	}
	res := make([]*api.Video, len(videos))
	for i, video := range videos {
		video.FavoriteCount = countList[i]
		apiVideo, _ := video.ToApiVideo(clientUser)
		res = append(res, apiVideo)
		apiVideo.IsFavorite = true
	}
	return res, nil
}
