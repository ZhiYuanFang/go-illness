package v1

import "github.com/gogf/gf/v2/frame/g"

type RefreshTokenReq struct {
	g.Meta       `path:"/refreshToken" method:"post"`
	RefreshToken string
}

type RefreshTokenRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
