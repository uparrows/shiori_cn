package model

import "errors"

var (
	ErrBookmarkNotFound  = errors.New("未找到书签")
	ErrBookmarkInvalidID = errors.New("无效的书签 ID")
	ErrTagNotFound       = errors.New("tag not found")

	ErrUnauthorized  = errors.New("未经授权的用户")
	ErrNotFound      = errors.New("未找到")
	ErrAlreadyExists = errors.New("已存在")
)
