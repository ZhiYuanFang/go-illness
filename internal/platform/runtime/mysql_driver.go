// Package runtime 注册进程级运行时依赖（对齐 go_ai_talk 手法，本仓仅 MySQL）。
package runtime

import (
	// 注册 GoFrame MySQL 驱动；须在任意 g.DB 之前导入。
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)
