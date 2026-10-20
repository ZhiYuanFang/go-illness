package service

import (
	"errors"
	"hello/internal/jwt"

	"github.com/gogf/gf/v2/net/ghttp"

)

type sMiddleware struct{}

var insMiddleware = sMiddleware{}

func Middleware() *sMiddleware {
	return &insMiddleware
}

func (s *sMiddleware) UserAuth(req *ghttp.Request) {
	//获取头部携带的token
	tokenString := req.Header.Get("Token")

	id, err := jwt.ParseToken(tokenString)
	//用户行动的用户ID
	if req.GetParam("id") == nil {
		req.SetParam("id", id)
	}

	if err == nil {
		if id > 0 {
			req.Middleware.Next()
		} else {
			req.SetError(errors.New("用户不存在"))
		}
	} else {
		req.SetError(err)
	}
}
