package utils

import "errors"

var (
	// DB errors
	ErrWrongToken             = errors.New("the token is wrong")
	ErrUserNotFound           = errors.New("the user is not found in database")
	ErrVideoUserNotExist      = errors.New("the author of the video is not found in database")
	ErrTokenVerifiedFailed    = errors.New("failed to verify the token in database")
	ErrGetFeedVideoListFailed = errors.New("failed to get feed video list")
	ErrGetUserVideoListFailed = errors.New("failed to get user's video list")
	ErrWrongParam             = errors.New("Wrong Parameter has been given")
	ErrTypeNotSupport         = errors.New("Type Not Support")
	ErrMathRealationFailed    = errors.New("failed to get relation of a user to a given token")
	ErrGetCommentListFailed   = errors.New("failed to get comment list of a video from database")
	ErrVideoNotExist          = errors.New("failed to get a video of a given id from database")
	ErrGetVideoFromUSer       = errors.New("failed to get a favorite video list of a user")
	ErrLikeFaile              = errors.New("Like failed")
	ErrDelCommentNotExist     = errors.New("Del comment is not exist")
	// NetWork errors
	ErrIpInitFailed = errors.New("failed to init IP")

	// IO erros
	ErrSaveVideoFaile   = errors.New("save video failed")
	ErrDBSaveVideoFaile = errors.New("save video in database failed")
)
