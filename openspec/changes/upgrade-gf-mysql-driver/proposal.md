## Why

本地 `ILLNESS_DB_LINK` 已成功注入（日志可见主机与库名），但登录查库报错：`cannot find database driver for specified database type ""`。根因是本仓锁定 **GoFrame v2.0.6**：`gdb.SetConfigGroup` **不会**从 `mysql:...` Link 解析 `Type`，导致驱动查找键为空。兄弟仓 `go_ai_talk` 已用 **gf v2.6.3** + `contrib/drivers/mysql/v2` 空白导入，同款 `dbcfg` 可工作。用户选定方案 **B**：对齐兄弟仓升级并注册驱动，而非仅在 ConfigNode 硬写 Type。

**对照说明**：已读 `openspec/project.md` 与空壳基线 `v0.0.1`（首版空壳，本次为 ADDED）。DB 连接串仍经 `ILLNESS_DB_LINK` 注入的约定不变；本变更修复「注入后仍无法建连」的运行时缺口。活跃 change `revive-illness-deploy-baseline` 的 `illness-db-env` 描述连接注入，本变更补充驱动/框架版本前提。

## What Changes

- 将 `github.com/gogf/gf/v2` 从 **v2.0.6** 升级到与兄弟仓一致的 **v2.6.3**（或同主线兼容补丁）
- 增加依赖 `github.com/gogf/gf/contrib/drivers/mysql/v2`，并在启动路径 **空白导入** 以注册 MySQL 驱动
- 保持现有 `dbcfg.ApplyGroupFromEnv` 语义；升级后 `SetConfigGroup` 会从 Link 解析 Type
- 修复因升级引起的编译/API 不兼容（若有），以 `go build` 通过为准
- **不**改上传、Compose、业务登录协议本身

## Capabilities

### New Capabilities

- `gf-mysql-runtime`：GoFrame 版本与 MySQL 驱动注册——保证经 env 注入的 `database.default` Link 可被 ORM 识别并建连

### Modified Capabilities

- （无基线 capability；若收版时与 `illness-db-env` 合并，视为对其「注入后可用」的运行时补强，本 change 以 ADDED 交付）

## Impact

- `go.mod` / `go.sum`：升级 gf、新增 contrib mysql 驱动
- 启动入口：`main.go` 或 `internal/platform` 下新增空白导入（对齐 `go_ai_talk` 的 `shared/runtime` 手法，本仓可更薄）
- 风险：gf 大版本跨度可能触及少量已弃用 API，需编译修复
- 验收：`.\hack\run-local.ps1` 后 `/login` 不再出现 type `""` 驱动错误
