package ctl

import "context"

type key int64

var userKey key

type UserInfo struct {
	UID      uint
	UserName string
}

func NewContext(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, userKey, value)
}

func GetFromContext(ctx context.Context) *UserInfo {
	v, ok := FromContext(ctx)
	if ok {
		return v
	}
	return nil
}

func FromContext(ctx context.Context) (*UserInfo, bool) {
	v, ok := ctx.Value(userKey).(*UserInfo)
	return v, ok
}
