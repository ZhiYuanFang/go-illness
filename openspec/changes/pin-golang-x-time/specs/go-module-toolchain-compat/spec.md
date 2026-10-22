## ADDED Requirements

### Requirement: x/time 与 Go 1.19 构建器兼容

在仓库声明的 Go 版本与 Dockerfile builder（当前为 **Go 1.19**）下，模块图 MUST NOT 选中要求高于该工具链的 `golang.org/x/time` 版本。本变更要求将 `golang.org/x/time` 固定为与兄弟仓一致的 **v0.6.0**（或其 `go` 指令 ≤ 1.19 的等价兼容版本），使 `go build` 与 Docker 多阶段构建中的 `go build` 可完成编译。

#### Scenario: 模块版本已钉死

- **WHEN** 审查 `go.mod` / `go.sum` 中的 `golang.org/x/time`
- **THEN** 解析到的版本为 `v0.6.0`（或设计明确允许的 ≤1.19 兼容版本），而非要求 Go 1.26 的 `v0.16.x`

#### Scenario: Docker builder 可编译

- **WHEN** 使用 `manifest/docker/Dockerfile`（`golang:1.19-alpine`）执行镜像构建中的 `go build`
- **THEN** MUST NOT 再出现 `golang.org/x/time/...: cannot compile Go 1.26 code`（或同等「依赖 go 指令高于 builder」）错误

### Requirement: 构建验收以编译通过为准

完成依赖调整后，本仓库 MUST 能在无 CGO 的 linux/amd64 目标下成功 `go build` 产出主二进制（与 Dockerfile 中的构建命令一致）。

#### Scenario: 本地 go build 成功

- **WHEN** 在已 tidy 的模块状态下于仓库根执行与 Docker 等价的 `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build`
- **THEN** 命令退出码为 0
