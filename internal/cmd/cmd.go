package cmd

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"

	"hello/internal/controller"
	"hello/internal/platform/ossclient"
	"hello/internal/service"
)

// Upload 将 multipart 文件直传 OSS，响应 CDN 绝对 URL（兼容字段 upload-file、header md5）。
func Upload(r *ghttp.Request) {
	ctx := r.Context()
	client, err := ossclient.New(ctx)
	if err != nil {
		r.Response.WriteStatusExit(500, err.Error())
		return
	}
	files := r.GetUploadFiles("upload-file")
	if len(files) == 0 || files[0] == nil {
		r.Response.WriteStatusExit(400, "缺少 upload-file")
		return
	}

	md5 := strings.TrimSpace(r.Header.Get("md5"))
	var objectName string
	if md5 != "" {
		objectName = md5 + ".png"
	} else {
		objectName = randomObjectName() + ".png"
	}
	objectKey := client.Config().ObjectKey(objectName)

	// 幂等：同 md5 已存在则直接返回 CDN URL
	if md5 != "" {
		exists, existErr := client.Exists(objectKey)
		if existErr != nil {
			r.Response.WriteStatusExit(500, existErr.Error())
			return
		}
		if exists {
			r.Response.WriteExit(client.Config().CDNURL(objectKey))
			return
		}
	}

	f, err := files[0].Open()
	if err != nil {
		r.Response.WriteStatusExit(400, err.Error())
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		r.Response.WriteStatusExit(400, err.Error())
		return
	}
	if len(data) == 0 {
		r.Response.WriteStatusExit(400, "空文件")
		return
	}

	cdnURL, err := client.PutBytes(objectKey, data, "image/png")
	if err != nil {
		r.Response.WriteStatusExit(500, err.Error())
		return
	}
	r.Response.WriteExit(cdnURL)
}

// randomObjectName 在无 md5 header 时生成唯一对象名片段。
func randomObjectName() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",

		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.BindHandler("/", func(r *ghttp.Request) {
				r.Response.Write("病变")
			})

			// WebSocket 房间推送已停用（对局迁本地引擎）；保留端点便于旧客户端探测
			s.BindHandler("/person", func(r *ghttp.Request) {
				r.Response.WriteStatusExit(410, "websocket room gone; use local engine")
			})
			//api
			s.Use(ghttp.MiddlewareHandlerResponse) //统一格式处理
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.POST("/upload", Upload)
				group.Bind(
					controller.Single,
					controller.Login,
					controller.RefreshToken,
				)
				s.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(
						service.Middleware().UserAuth,
					)
					group.Bind(
						//提供用户信息查询接口
						controller.User,
						controller.Test,
						controller.NextIllness,
					)
				})
			})
			s.SetServerRoot(gfile.MainPkgPath()) //设置exe文件生成位子
			s.Run()
			return nil
		},
	}
)
