package utils

import "errors"

var (
	// DB errors
	ErrWrongToken   = errors.New("the token is wrong")
	ErrUserNotFound = errors.New("the user is not found in database")
)
