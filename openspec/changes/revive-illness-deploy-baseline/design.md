## Context

`go-illness` 为 GoFrame 单仓老代码：明文 DB 密码、Mac 本地上传目录、旧 Dockerfile（`ADD ./main`）、无 Compose/ACR/.env。兄弟仓 `go_ai_talk` 提供可参考的部署/配置手法，但本仓约束为：**单进程、无 Redis/MQ、仓内无 Nginx、上传 API 本变更不动**。

当前 OpenSpec 基线 `v0.0.1` 为空壳；另有已完成未收版的 `align-openspec-versioned-baseline`。本设计只覆盖部署与配置基线。

## Goals / Non-Goals

**Goals:**

- 可用 Compose 拉起唯一 `illness-service`（端口 9686）
- CI tag 构建推送 ACR；`IMAGE_TAG` 人工控制（起步 `v0.0.1`）
- DB：`ILLNESS_DB_LINK` → 同实例库 `illness`；yaml 无真实密钥
- OSS 配置约定：同桶同 CDN，`illness/` 前缀；密钥走 env（供后续上传 change 使用）
- 更新 `project.md` 写入上述全局约定

**Non-Goals:**

- 不改 `POST /upload` 及静态 `/images` 行为
- 不引入 Redis/MQ/Nginx 文件/上传 volume
- 不拆微服务、不大重构业务包

## Decisions

### D1：单服务 Compose，对齐兄弟仓「基线 + overlay」缩小版

- **选择**：`docker-compose.yml`（本地 build / `:local`）+ `docker-compose.test.yml` / `docker-compose.prod.yml`（`REGISTRY`/`IMAGE_TAG` pull）；`.env.example` 为唯一提交的敏感信息模板。
- **理由**：与运维习惯一致且本仓只有一个进程。
- **备选**：单文件 compose —— 拒绝，test/prod 镜像源差异需要 overlay。

### D2：ACR 工作流单 matrix

- **选择**：从 `go_ai_talk` 的 `docker-acr.yml` 收缩为单一镜像（如 `illness-service`）；触发仍为 `v*` tag + `workflow_dispatch`；prod/test 由 tag 形态区分（纯 `vX.Y.Z` → prod，含预发布后缀 → test）。
- **理由**：版本由人工打 tag / 改 `.env` 的 `IMAGE_TAG`，与用户约定一致。
- **备选**：仅 manual dispatch —— 可保留，但不替代 tag 主路径。

### D3：DB 注入用薄封装，不整包抄 platform

- **选择**：启动时读取 `ILLNESS_DB_LINK`（兼容 `GF_DATABASE_DEFAULT_LINK`），写入 GoFrame `database.default.link`；`config.yaml` 中 `link` 置空或占位注释。可最小实现于 `main`/`cmd` 或极薄 `internal/platform/dbcfg`。
- **理由**：手法同源即可，避免无关 platform 依赖。
- **备选**：整目录复制 `go_ai_talk/internal/platform` —— 拒绝。

### D4：OSS 本变更只落配置与约定

- **选择**：`config.yaml`（或专用段）写入与兄弟仓一致的桶/region/cdn 占位，`objectKeyPrefix: illness/`；`.env.example` 增加 `ILLNESS_OSS_ACCESS_KEY_ID` / `SECRET`（及可选 `ILLNESS_OSS_ENDPOINT`）；**不**改 Upload handler。
- **桶/CDN**：与 `go_ai_talk` 同（`caidoukeji` + `https://resorce.cuplay.top`，以实现时核对兄弟仓当前值为准）。
- **理由**：用户确认同桶同 CDN + `illness/`，且上传接口后续再改。
- **备选**：本变更直接改上传 —— 明确排除。

### D5：端口与库名

- **选择**：`server.address: ":9686"`；MySQL database name = `illness`。
- **理由**：用户已确认。

### D6：明文密钥与旧上传路径

- **选择**：移除 yaml 中真实密码；`.gitignore` 确保 `.env` / `.env.local` 等不入库。旧 `FileDir` 上传逻辑本变更保留（已知容器内不可用），在 design/tasks 标注为后续 OSS 变更债务，不在本 change 修接口。

## Risks / Trade-offs

- [本变更后容器内上传仍依赖 Mac 路径] → Mitigation：文档/tasks 标明；后续 OSS change 删除本地盘逻辑；部署验收不把上传列为必测。
- [历史明文密码曾进 git] → Mitigation：轮换 DB 密码；仅 `.env` 持有新密钥。
- [ACR 命名空间/Secrets 未建] → Mitigation：tasks 列出运维清单；workflow 可先合入，首次 tag 前配 Secrets。
- [与未收版 OpenSpec 工程 change 并存] → Mitigation：两 change 独立；收版版本号由人工决定（建议工程约定与部署基线可同收 `v0.0.2` 或分版）。

## Migration Plan

1. 合入本变更 → 配置 ACR Secrets / 创建库 `illness` → 复制 `.env.example` 为服务器 `.env`，填 `ILLNESS_DB_LINK`、`IMAGE_TAG=v0.0.1`。
2. `git tag v0.0.1 && push` → Actions 推镜像 → compose overlay pull up。
3. 回滚：改回上一 `IMAGE_TAG` 或本地旧二进制；DB 仅连接串变更无 schema 迁移（本变更不改表）。

## Open Questions

- 无阻塞项。ACR 具体 namespace 名以实现时与运维确认写入 `.env.example` 注释即可。
