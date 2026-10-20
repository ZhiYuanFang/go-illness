package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type UserInfoReq struct {
	g.Meta `path:"/UserInfo" method:"get"`
	Id     int64
}

type UserInfoRes struct {
	Id                int64       `json:"id"`       // User ID
	PassWord          string      `json:"passWord"` // User Password
	NickName          string      `json:"nickName"` // User Nickname
	Avatar            string      `json:"avatar"`   //头像
	Lv                float32     `json:"lv"`       // 等级
	CreateAt          *gtime.Time `json:"createAt"` // Created Time
	UpdateAt          *gtime.Time `json:"updateAt"` // Updated Time
	OpenId            string      `json:"openId"`
	Rank              int         `json:"rank"`
	MaxRank           int         `json:"maxRank"`
	UseLiveDay        *gtime.Time `json:"useLiveDay"`
	TutorialCompleted bool        `json:"tutorialCompleted"` // 教学是否已完成
}

type UserInfoUpdateReq struct {
	g.Meta   `path:"/UserInfoUpdate" method:"post"`
	Id       int64
	Avatar   string
	NickName string
}

type UserRankUpdateReq struct {
	g.Meta `path:"/UserRankUpdate" method:"post"`
	Id     int64
	Rank   int
}

type UserReliveUpdateReq struct {
	g.Meta `path:"/reLive" method:"post"`
	Id     int64
}

type UserReliveJudgeReq struct {
	g.Meta `path:"/reLiveJudge" method:"get"`
	Id     int64
}
