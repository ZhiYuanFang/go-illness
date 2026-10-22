## Why

本地/构建机拉取基础镜像失败：`m.daocloud.io` DNS 解析 `no such host`。`Dockerfile` 将 Docker Hub 官方镜像硬编码为 DaoCloud 代理前缀，域名失效或环境 DNS 不认时整条构建链断裂。用户选定方案 **A**：Dockerfile 使用标准官方镜像名，镜像加速交由本机/CI 的 registry mirror 配置，不把某一家代理写死进仓库。

**对照说明**：已读 `openspec/project.md` 与空壳基线 `v0.0.1`（首版空壳，本次为 ADDED）。单服务 Compose/ACR 发版约定不变；仅修正基础镜像引用方式。

## What Changes

- `manifest/docker/Dockerfile` 的 builder / runtime `FROM` 改为官方名：`golang:1.19-alpine`、`alpine:3.19`（去掉 `m.daocloud.io/docker.io/library/` 前缀）
- 保留 Dockerfile 内 `GOPROXY`（go module 加速，与 registry 无关）
- 不在仓库内绑定新的第三方 registry 前缀；国内加速说明可简短写在 README 或 Dockerfile 注释（可选）
- 不改业务代码、不改 ACR workflow 推送逻辑（除非文档需同步一句）

## Capabilities

### New Capabilities

- `docker-base-images`：illness-service Dockerfile 基础镜像引用规范——使用官方镜像名，registry 加速由构建环境配置

### Modified Capabilities

- （无基线 capability；与未收版 `single-service-deploy` 相关但属实现细节修正，本变更以 ADDED 交付）

## Impact

- 文件：`manifest/docker/Dockerfile`（必要时 README 一句说明）
- 本地：需能直连 Docker Hub，或在 Docker daemon 配置 mirror
- CI（GitHub Actions `ubuntu-latest`）：通常可直连官方镜像；不再依赖 DaoCloud 域名
