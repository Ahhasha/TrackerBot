package model

import "errors"

var (
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInvalidCategory   = errors.New("invalid category")
	ErrInvalidUser       = errors.New("invalid user")
	ErrUserNotRegistered = errors.New("user not registered")
	ErrCategoryNotFound  = errors.New("category not found")
)
