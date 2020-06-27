package internal

import "errors"

// Generic application errors
var (
	ErrNotFound     = errors.New("not found")
	ErrTokenExpired = errors.New("token expired")
)
