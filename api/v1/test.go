package v1

import "github.com/gogf/gf/v2/frame/g"

type TestReq struct {
	g.Meta `path:"/test" method:"get"`
	Id     int
}
