## ADDED Requirements

### Requirement: Dockerfile 使用官方基础镜像名

`manifest/docker/Dockerfile` 的多阶段构建 MUST 使用 Docker Hub 官方镜像短名作为 `FROM` 源，至少包括 builder 阶段的 **`golang:1.19-alpine`**（或与 `go.mod` 声明的 Go 主版本一致的官方 golang alpine 标签）与 runtime 阶段的 **`alpine:3.19`**（或设计锁定的等价 alpine 小版本）。MUST NOT 将某一第三方 registry 代理前缀（例如 `m.daocloud.io/docker.io/library/`）硬编码进仓库 Dockerfile。

#### Scenario: 无 DaoCloud 硬编码前缀

- **WHEN** 审查 `manifest/docker/Dockerfile` 的全部 `FROM` 行
- **THEN** 不存在 `m.daocloud.io`（或同类写死代理主机）前缀；基础镜像名为官方短名

#### Scenario: 构建可解析官方标签

- **WHEN** 构建环境能够解析并拉取 Docker Hub 官方镜像（或通过 daemon registry-mirrors 透明加速）
- **THEN** 使用本 Dockerfile 的 `docker build` / Compose `--build` 不再因 `m.daocloud.io` DNS 失败而在拉基础镜像阶段中断

### Requirement: go module 代理与 registry 分离

Dockerfile 内 MAY 继续配置 `GOPROXY` 以加速 `go mod download`。该配置 MUST NOT 被视为基础镜像 registry 的替代方案。

#### Scenario: GOPROXY 可保留

- **WHEN** 审查 Dockerfile 构建阶段环境变量
- **THEN** 允许存在 `GOPROXY=...goproxy.cn...` 一类设置，且与官方 `FROM` 镜像名同时存在
