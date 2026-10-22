## Why

Docker 构建（`golang:1.19-alpine`）失败：`golang.org/x/time/rate: cannot compile Go 1.26 code`。当前 `go.mod` 间接依赖 `golang.org/x/time v0.16.0`（模块声明 `go 1.26.0`），与 Dockerfile / `go.mod` 的 Go 1.19 不兼容。兄弟仓 `go_ai_talk` 将 `x/time` 钉在 **v0.6.0**。用户选定：**改 `x/time` 引用（降版本钉死）**，不在本变更抬高构建器到 1.26。

**对照说明**：已读 `openspec/project.md` 与空壳基线 `v0.0.1`（首版空壳，本次为 ADDED）。官方 `FROM` 策略见已完成的 `fix-docker-base-image-mirror`；本变更只动模块版本。

## What Changes

- 将 `golang.org/x/time` 固定为与 Go 1.19 兼容的版本（目标 **v0.6.0**，对齐 `go_ai_talk`）
- `go mod tidy`；必要时对仍要求过高 Go 版本的其它 `golang.org/x/*` 做同类钉版本（仅在构建仍失败时）
- **不**修改 Dockerfile Go 小版本（保持 `golang:1.19-alpine`），除非钉依赖后仍无法通过再另开变更

## Capabilities

### New Capabilities

- `go-module-toolchain-compat`：Docker/本地 Go 1.19 工具链与间接依赖 `go` 指令兼容——禁止拉入要求高于构建 Go 版本的 `golang.org/x/time`（及同类已识别阻塞项）

### Modified Capabilities

- （无）

## Impact

- `go.mod` / `go.sum`
- Docker `go build` 与本地 `go build` 验收
- 间接来自 OSS SDK 等依赖的版本选择；功能行为预期不变
