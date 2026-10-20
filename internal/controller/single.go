package controller

import (
	"context"
	"fmt"
	v1 "hello/api/v1"
)

type cSingle struct{}

var Single = cSingle{}

var singleID = -1

func (c *cSingle) GetSingle(ctx context.Context, req *v1.SingleReq) (res *v1.SingleRes, err error) {
	singleID = singleID - 1
	fmt.Println("singleId", int64(singleID))
	return &v1.SingleRes{
		SingleID: int64(singleID),
	}, nil
}
