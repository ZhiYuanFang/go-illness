## 1. Dockerfile

- [x] 1.1 将 builder `FROM` 改为 `golang:1.19-alpine`，runtime `FROM` 改为 `alpine:3.19`；删除 `m.daocloud.io/docker.io/library/` 前缀
- [x] 1.2 保留 `GOPROXY`；可在 Dockerfile 顶部用简短中文注释说明：registry 加速请配本机/CI mirror，勿写死代理前缀

## 2. 文档（可选但推荐）

- [x] 2.1 在 `README.MD` 部署/构建相关处补一句：基础镜像为官方名；国内需自配 Docker registry-mirrors

## 3. 校验

- [x] 3.1 `openspec validate fix-docker-base-image-mirror --strict`
