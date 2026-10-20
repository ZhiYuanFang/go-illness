package v1

import (
	entity "hello/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type RankListReq struct {
	g.Meta `path:"/RankList" method:"get"`
}

type RankListRes struct {
	UserList []*entity.User
}
