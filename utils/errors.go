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

	// NetWork errors
	ErrIpInitFailed = errors.New("failed to init IP")

	// IO erros
	ErrSaveVideoFaile   = errors.New("save video failed")
	ErrDBSaveVideoFaile = errors.New("save video in database failed")
)
