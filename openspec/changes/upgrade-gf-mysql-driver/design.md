## Context

`dbcfg` 仅设置 `ConfigNode{Link, Role}`。gf **v2.6.3** 的 `SetConfigGroup` 会 `parseConfigNode` 从 Link 填 Type；**v2.0.6** 不会。新版 MySQL 驱动已迁出 core，须空白导入 `contrib/drivers/mysql/v2`。本仓 go 版本已为 1.19，与兄弟仓一致。

## Goals / Non-Goals

**Goals:**

- gf 升级到 **v2.6.3**（与 `go_ai_talk` 对齐）
- 注册 MySQL contrib 驱动
- 登录等首次 `g.DB()` 访问可建连（Type 非空且驱动已注册）
- `go build` 通过

**Non-Goals:**

- 不引入 Redis contrib（本仓无 Redis）
- 不改 `ILLNESS_DB_LINK` / Compose 注入约定
- 不做方案 A 单独硬编码 Type（升级后由框架解析；若编译期仍需可兼写 Type 作双保险，非必须）
- 不大规模重构业务代码，仅修升级引起的编译错误

## Decisions

### D1：目标版本锁定 v2.6.3

- **选择**：`github.com/gogf/gf/v2 v2.6.3` + `github.com/gogf/gf/contrib/drivers/mysql/v2 v2.6.3`
- **理由**：与兄弟仓一致，减少「两边行为不一致」；用户选定 B。

### D2：空白导入位置

- **选择**：在 `main.go`（或极薄 `internal/platform/runtime` 包）空白导入 `_ "github.com/gogf/gf/contrib/drivers/mysql/v2"`，保证任意 `g.DB` 之前已注册。
- **理由**：对齐 go_ai_talk `internal/shared/runtime`；本仓单服务，不必复制 Redis 导入。

### D3：dbcfg 是否兼写 Type

- **选择**：优先依赖升级后的 `parseConfigNode`；若实现时发现边界情况，可额外从 Link 前缀设置 `Type`（防御性，非目标）。
- **理由**：B 的核心是框架对齐，而非绕过。

### D4：升级策略

- **选择**：`go get` 指定版本 → `go mod tidy` → `go build` → 按编译错误最小修补。
- **理由**：跨 2.0→2.6 可能有 API 微调；禁止无关重构。

## Risks / Trade-offs

- [API 破坏面] → Mitigation：只修编译失败点；对照 go_ai_talk 同类用法
- [间接依赖膨胀] → Mitigation：`go mod tidy`；不主动加 Redis 等无关 contrib
- [Docker 镜像构建需拉新模块] → Mitigation：常规；CI/本地 GOPROXY 沿用 goproxy.cn

## Migration Plan

1. 升级依赖并加空白导入 → 本地 `go build` / `.\hack\run-local.ps1`
2. 小游戏再打 `/login`：应越过「type ""」错误（后续若缺表/缺库再另排）
3. 回滚：恢复 go.mod 旧版本（不推荐；登录仍坏）

## Open Questions

- （无阻塞）
