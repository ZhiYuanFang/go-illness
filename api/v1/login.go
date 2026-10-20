package v1

import "github.com/gogf/gf/v2/frame/g"

type LoginReq struct {
	g.Meta `path:"/login" method:"post"`
	WXCode string `json:"wxCode"`
}

type LoginRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
