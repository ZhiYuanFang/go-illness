## 1. 配置与密钥清理

- [x] 1.1 清理 `manifest/config/config.yaml`：移除真实 DB 密码；保留 `:9686`；增加 OSS 段占位（同桶同 CDN、`objectKeyPrefix: illness/`，密钥留空）
- [x] 1.2 实现启动时 `ILLNESS_DB_LINK`（兼容 `GF_DATABASE_DEFAULT_LINK`）注入 `database.default.link`
- [x] 1.3 更新 `.gitignore`，确保 `.env`、`.env.local`、`.env.prod`、`.env.test` 等不被提交

## 2. Docker Compose 与镜像

- [x] 2.1 重写 `manifest/docker/Dockerfile` 为可在 CI 从源码构建的多阶段镜像（单服务）
- [x] 2.2 新增 compose 基线 + test/prod overlay（仅 illness-service；无 redis/mq/nginx；无上传 volume）
- [x] 2.3 新增 `manifest/docker/.env.example`（`REGISTRY`、`IMAGE_TAG=v0.0.1`、`ILLNESS_DB_LINK`、OSS AK/SK 等占位）

## 3. ACR 工作流

- [x] 3.1 新增 `.github/workflows/docker-acr.yml`（单服务 matrix；`v*` tag + workflow_dispatch；对齐 go_ai_talk 的 test/prod 解析思路）
- [x] 3.2 在 `.env.example` 或简短注释中写明需配置的 GitHub Secrets / Environment（`ACR_USERNAME`、`ACR_PASSWORD`、`REGISTRY`）

## 4. 文档约定与校验

- [x] 4.1 更新 `openspec/project.md`：补充单服务部署、DB env、OSS 同桶同 CDN + `illness/`、无 Redis/MQ、仓内无 Nginx、上传改造另开 change 等约定
- [x] 4.2 确认未改动现有上传接口代码路径/契约；`openspec validate revive-illness-deploy-baseline --strict`
