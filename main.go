package main

import (
	"hello/internal/cmd"
	"hello/internal/platform/dbcfg"
	_ "hello/internal/platform/runtime"

	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	// 数据库连接串经环境变量注入（Compose/.env）；yaml 不落真实密钥。
	dbcfg.ApplyGroupFromEnv("illness-service", "default", "ILLNESS_DB_LINK", "GF_DATABASE_DEFAULT_LINK")
	cmd.Main.Run(gctx.New())
}
