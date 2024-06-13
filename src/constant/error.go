package constant

import "errors"

var (
	ErrExpired      = errors.New("expired")
	ErrRegistered   = errors.New("registered")
	ErrRequired     = errors.New("required")
	ErrInvalid      = errors.New("invalid")
	ErrInvalidRange = errors.New("invalid range")
	ErrNotFound     = errors.New("not found")
	ErrCompleted    = errors.New("completed")
)
