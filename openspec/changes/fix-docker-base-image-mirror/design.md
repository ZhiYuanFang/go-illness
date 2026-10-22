## Context

当前 Dockerfile：

```
FROM m.daocloud.io/docker.io/library/golang:1.19-alpine AS builder
FROM m.daocloud.io/docker.io/library/alpine:3.19
```

构建时报 `lookup m.daocloud.io … no such host`。GOPROXY 只影响 `go mod`，不解决基础镜像拉取。

## Goals / Non-Goals

**Goals:**

- `FROM` 使用官方短名 `golang:1.19-alpine`、`alpine:3.19`
- 构建环境可自行配置 registry-mirrors；仓库不绑定单一国内代理

**Non-Goals:**

- 不升级 Go 基础镜像大版本（仍 1.19，与 go.mod 一致）
- 不引入 build-arg 切换 mirror（除非后续明确需要）
- 不改 ACR 推送与 Compose 拓扑

## Decisions

### D1：官方镜像名（方案 A，已锁定）

- **选择**：去掉一切 `m.daocloud.io/...` 前缀。
- **理由**：代理域名不可控；GitHub Actions 与多数云构建默认可拉 Hub。
- **备选**：换另一家前缀——已否决（仍写死）。

### D2：加速责任边界

- **选择**：国内慢速时由运维在 Docker Desktop / daemon.json / 构建机配置 `registry-mirrors`。
- **理由**：与「密钥走 env、不进库」同理——环境相关配置不硬编码进 Dockerfile。

### D3：GOPROXY 保留

- **选择**：继续 `ENV GOPROXY=https://goproxy.cn,...`
- **理由**：与 registry 无关，对国内 `go mod download` 仍有价值。

## Risks / Trade-offs

- [国内直连 Hub 慢或超时] → Mitigation：本机配 mirror；CI 一般可用；必要时再开 build-arg 变更
- [标签拉取策略差异] → Mitigation：固定 `1.19-alpine` / `3.19`，与现意图一致

## Migration Plan

1. 合入 Dockerfile → 本地 `docker compose … --build` 或 `docker build -f manifest/docker/Dockerfile`
2. 若仍失败，检查本机能否 `docker pull golang:1.19-alpine`，再配 mirror

## Open Questions

- （无阻塞）
