package controller

import (
	"context"
	v1 "hello/api/v1"
)

type cRank struct{}

var Rank = cRank{}

func (c *cRank) UpDateUserRank(ctx context.Context, req *v1.RankListReq) (res *v1.RankListRes, err error) {
	return nil, nil
}
