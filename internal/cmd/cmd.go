package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"

	"hello/internal/consts"
	"hello/internal/controller"
	"hello/internal/service"
)

// MAC 本地服务器需要
var FileDir = "/Users/mac/Desktop/go_server_illness_images/"
var StaticPath = "images"

func Upload(r *ghttp.Request) {
	files := r.GetUploadFiles("upload-file") //图片校验交给客户端重命名md5值
	md5 := r.Header.Get("md5")
	md5Png := md5 + ".png"
	if gfile.Exists(FileDir + md5Png) {
		println("文件已存在")
		r.Response.WriteExit(StaticPath + "/" + md5Png)
		return
	} else {
		//修改files的文件名
		files[0].Filename = md5Png
		println("存入的文件")
		println(files[0].Filename)
		_, err := files.Save(FileDir)
		if err != nil {
			r.Response.WriteExit(err)
		} else {
			r.Response.WriteExit(StaticPath + "/" + files[0].Filename)
		}

	}
}

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",

		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 如果FileDir不存在则创建
			if !gfile.Exists(FileDir) {
				gfile.Mkdir(FileDir)
			}
			s := g.Server()
			s.SetFileServerEnabled(true)
			s.SetIndexFolder(true)
			s.AddStaticPath("/"+StaticPath, FileDir)
			s.BindHandler("/", func(r *ghttp.Request) {
				r.Response.Write("病变")
			})

			// ******** 链接websocket的时候 携带用户ID、房间ID ws://127.0.0.1:9686/person?id=1
			s.BindHandler("/person", func(r *ghttp.Request) {
				var ctx = r.Context()
				ws, err := r.WebSocket()
				playerid, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
				if err != nil {
					fmt.Println("ID错误")
					r.Exit()
				}
				consts.PersonWSArr[playerid] = ws
				fmt.Println(fmt.Sprintf("有用户进入长连接，当前长连接的集合结果：", consts.PersonWSArr))
				if err != nil {
					fmt.Println(r.URL)
					glog.Error(ctx, err)
					r.Exit()
				}
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
			// s.EnableHTTPS("/Users/mac/Desktop/lovelin/lovelin.cuplay.top.pem", "/Users/mac/Desktop/lovelin/lovelin.cuplay.top.key")
			s.Run()
			return nil
		},
	}
)
