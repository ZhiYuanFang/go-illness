## Why

本仓业务代码处于长期弃用、需重新捡起开发的状态：配置含明文密钥、Docker 依赖本机预编译、无 Compose/ACR/.env 闭环，无法按兄弟仓 `go_ai_talk` 的方式稳定部署。需要先落地「单进程部署基线」与数据库/OSS 配置约定，为后续业务迭代提供可运行骨架；本变更明确不拆微服务、不引入 Redis/MQ。

## What Changes

- 引入单服务 docker-compose（基线 + test/prod overlay）、现代化多阶段 Dockerfile、`.env.example`（敏感信息仅经环境变量）
- 新增 `.github/workflows/docker-acr.yml`：tag 触发构建并推送阿里云 ACR；镜像主 tag 与服务器 `.env` 的 `IMAGE_TAG` 对齐，起步示例为 `v0.0.1`，后续版本由人工控制
- 数据库配置手法对齐兄弟仓：yaml 去掉明文 `link`；进程通过 `ILLNESS_DB_LINK`（及必要时 `MYSQL_TCP_HOST`）注入；库名 **`illness`**，与 `go_ai_talk` **同 MySQL 实例、不同库**；HTTP 端口保持 **`:9686`**
- 约定 OSS：**与 go_ai_talk 同桶同 CDN**，`objectKey` 前缀 **`illness/`**；密钥与 endpoint 经 `.env`/Compose 注入；**本次不改上传接口实现与路由**
- **不**在仓库内新增 Nginx 配置；**不**引入 Redis/MQ；**不**为上传目录做本地 volume 挂载（上传改造留待后续 change）
- 同步更新 `openspec/project.md` 中与部署/DB/OSS 相关的全局约定（实现阶段）

### 非目标（Non-goals）

- 不修改现有 `POST /upload` 等上传接口行为与响应格式（后续另开 change 改为 OSS 直传并返回 CDN）
- 不拆微服务、不迁移 `internal/services` 大重构
- 不升级 Go/GoFrame 大版本（除非构建所必需的最小调整）
- 不收版 OpenSpec（本 change 与已完成的 `align-openspec-versioned-baseline` 可并存，由人工决定何时 `/opsx-archive`）

## Capabilities

### New Capabilities

- `single-service-deploy`: 单进程 docker-compose、Dockerfile、ACR 工作流、`.env` 敏感信息与 `IMAGE_TAG` 发版约定
- `illness-db-env`: MySQL `illness` 库经 `ILLNESS_DB_LINK` 注入、配置去密钥、监听端口 `:9686`
- `oss-cdn-illness-prefix`: OSS/CDN 配置基线（同桶同 CDN、前缀 `illness/`、密钥走环境变量）；不含上传 API 改造

### Modified Capabilities

- （无）当前基线 `openspec/specs/v0.0.1` 为空壳，本变更为首批 ADDED capability

## Impact

- 影响：`manifest/docker/**`、`.github/workflows/**`、`manifest/config/config.yaml`、启动入口（env 覆盖 DB）、`.gitignore`（若需忽略 `.env*`）、`openspec/project.md`
- 不影响：现有业务路由语义、上传接口契约、客户端协议（除运维侧改用新镜像与 `.env`）
- 依赖：阿里云 ACR 命名空间与 GitHub Secrets（`ACR_USERNAME`/`ACR_PASSWORD`/`REGISTRY`）；同实例上已存在或可创建的 MySQL 库 `illness`
- 基线依据：空壳 `v0.0.1` + 本 change delta（首版空壳，本次为 ADDED）
