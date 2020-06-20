package internal

import "errors"

var (
	// ErrNotFound signals that a resource wasn't found
	ErrNotFound = errors.New("not found")
)
