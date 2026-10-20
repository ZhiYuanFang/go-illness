package controller

import (
	"context"
	"errors"
	"fmt"
	v1 "hello/api/v1"
)

type cTest struct{}

var Test = cTest{}

func (t *cTest) Test(ctx context.Context, req *v1.TestReq) (res *v1.LoginRes, err error) {
	if req.Id > 0 {
		err = errors.New(fmt.Sprintf("%s%d", "用户ID：", req.Id))
	} else {
		err = errors.New("用户不存在")
	}
	return nil, err
}
