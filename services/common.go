package services

import "errors"

var (
	ErrUnexpectedDBError = errors.New("unexpected database error")
	ErrUserNotAuthorized = errors.New("user is not authorized")
)
