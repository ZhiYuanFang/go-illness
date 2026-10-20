package v1

import "github.com/gogf/gf/v2/frame/g"

type SingleReq struct {
	g.Meta `path:"/getSingleId" method:"get"`
}
type SingleRes struct {
	SingleID int64 `json:"singleId"`
}
