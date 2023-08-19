package utils

import "errors"

var (
	// DB errors
	ErrWrongToken        = errors.New("the token is wrong")
	ErrUserNotFound      = errors.New("the user is not found in database")
	ErrVideoUserNotExist = errors.New("the author of the video is not found in database")

	// NetWork errors
	ErrIpInitFailed = errors.New("failed to init IP")
)
