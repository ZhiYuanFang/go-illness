## Context

`aliyun-oss-go-sdk` 等依赖通过 `go get`/`tidy` 将 `golang.org/x/time` 拉到 v0.16.0；该版本 `go.mod` 为 `go 1.26.0`。构建镜像为 Go 1.19，故失败。`go_ai_talk` 使用 `x/time v0.6.0`（`go 1.18`）。

## Goals / Non-Goals

**Goals:**

- 钉 `golang.org/x/time@v0.6.0`
- `go build`（本机与 Docker builder 1.19）通过

**Non-Goals:**

- 不升级 Dockerfile 到 Go 1.23/1.26（可后续对齐 cash-service）
- 不更换 OSS SDK 大版本（除非钉 time 后无法满足最小依赖，再评估）

## Decisions

### D1：钉 v0.6.0（已锁定）

- **选择**：与兄弟仓相同的 `v0.6.0`。
- **理由**：已知与 Go 1.19 兼容；减少两边差异。
- **备选**：升构建器到 1.26——用户已否，优先改引用。

### D2：实施命令

- **选择**：`go get golang.org/x/time@v0.6.0` → `go mod tidy` → `go build`；Docker 再编一次确认。
- **理由**：最小路径。

### D3：连锁依赖

- **选择**：若 tidy 后仍有其它模块报「cannot compile Go 1.2x」，按同样原则钉到兄弟仓或 Go 1.19 可接受版本；记入 tasks 实现备注。
- **理由**：一次修通构建，避免只换 time 又卡下一包。

## Risks / Trade-offs

- [OSS SDK 声明需要更新的 x/time] → Mitigation：v0.6 通常仍满足 rate 限流 API；以编译与上传冒烟为准
- [日后 go get 再次漂回 v0.16] → Mitigation：可在 go.mod 保留显式 require；或后续升 Go 工具链

## Migration Plan

1. 钉版本并 tidy → 本地 build → Docker build
2. 无需数据迁移

## Open Questions

- （无阻塞）
