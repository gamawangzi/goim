package comet

import "errors"

var (
	ErrUserNotOnline = errors.New("用户不在线")
	ErrAuthFailed    = errors.New("认证失败")
)
