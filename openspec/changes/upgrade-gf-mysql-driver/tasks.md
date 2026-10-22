## 1. 依赖升级

- [x] 1.1 将 `github.com/gogf/gf/v2` 升级到 `v2.6.3`，并添加 `github.com/gogf/gf/contrib/drivers/mysql/v2@v2.6.3`
- [x] 1.2 `go mod tidy`；确认 `go.mod`/`go.sum` 版本一致

## 2. 驱动注册与编译

- [x] 2.1 在 `main.go` 或 `internal/platform/runtime` 空白导入 MySQL contrib 驱动（中文注释说明用途）
- [x] 2.2 `go build` 通过；若有 API 不兼容，做最小修补（禁止无关重构）

## 3. 校验

- [x] 3.1 `openspec validate upgrade-gf-mysql-driver --strict`
- [x] 3.2 （手工）`.\hack\run-local.ps1` 后登录：确认不再出现 type `""` 驱动错误
