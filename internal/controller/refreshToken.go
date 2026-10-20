package controller

import (
	"context"
	v1 "hello/api/v1"
	"hello/internal/jwt"
)

type cRefreshToken struct{}

var RefreshToken = cRefreshToken{}

func (r *cRefreshToken) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
	id, err := jwt.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	//refreshToken 有效
	//执行token获取操作
	accessToken, refreshToken, err := jwt.NewToken(id)
	if err != nil {
		return nil, err
	}
	//返回前端数据
	return &v1.RefreshTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, err
}
