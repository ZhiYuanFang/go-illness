# 项目上下文

## 项目目标

本项目（go-illness）基于 GoFrame 构建业务服务能力，并以 OpenSpec 管理变更提案、设计、任务与版本化规格基线。当前通过版本化 OpenSpec 工程约定，保证每次新建需求都能对照历史规格评估对全局代码的影响。

## 技术栈

- Go（Golang）
- GoFrame（HTTP 服务、路由、中间件、配置）
- MySQL（库名 `illness`；与兄弟仓 `go_ai_talk` 可同实例、不同库）
- 阿里云 OSS + CDN（与 go_ai_talk 同桶同 CDN，objectKey 前缀 `illness/`；上传接口改造另开 change）
- Docker Compose + 阿里云 ACR（单服务发版）
- OpenSpec（变更提案、设计、任务与规格管理）

## 项目约定

### 文档语言约定（强制）

- 所有由 OpenSpec 生成的文档（proposal、design、specs、tasks）必须使用中文撰写。
- 面向团队协作的说明、约束、验收标准、任务拆解均使用中文表述。
- 仅在必须保留英文标识时使用英文（例如环境变量名、接口路径、技术关键字）。

### 代码注释约定（强制）

- 后续生成或修改代码时，必须尽可能补充中文注释，重点覆盖：复杂流程、跨模块调用、重试/降级/错误语义、关键配置项。
- 公共导出函数、结构体与关键字段优先添加中文注释。
- 允许保留英文符号名、库 API 与协议字段名，但解释说明必须使用中文。
- 禁止添加“废话注释”；注释应解释“为什么这样做”与“边界条件”。

### 测试文件约定（强制）

- 在后续代码生成与实现过程中，不允许新建任何测试文件。
- 禁止创建或新增以下文件类型：`*_test.go`、`*.spec.*`、`*.test.*`。
- 若历史测试文件存在，按当前变更策略逐步清理，不作为新实现的交付前置条件。

### OpenSpec 基线参考约定（强制）

- 后续 AI 发起任何新变更前，必须先读取并对照 **`openspec/specs/v0.0.1/spec.md`**（当前合并基线规格），再生成 proposal/design/tasks。
- 若当前基线为空壳（无 capability），仍必须阅读本文件（`project.md`）与该空壳基线以确认指针与约定；不得跳过。
- 若本次需求涉及已有 capability，必须在 proposal 中明确标注“复用/变更了哪些已有 spec 边界”；禁止脱离历史 spec 直接重写同类能力。
- AI 在实现阶段（apply）必须以当前基线中的 Requirement/Scenario 作为验收约束（基线有内容时）；若实现与基线冲突，需先更新变更规格再改代码。
- AI 在评审或总结阶段必须说明“本次改动对应的 spec 依据”，至少列出相关 capability 名称（基线目录中的章节名；空壳期可写明「首版空壳，本次为 ADDED」）。
- 新版本发版时可运行 `hack/archive-openspec-to-version.py vX.Y.Z`（或等价流程）合并为 `openspec/specs/vX.Y.Z/spec.md`；若行为发生变更而未同步到 specs，不允许归档为完成状态（除非明确记录 `--skip-specs` 原因并经人工确认）。
- **收版语义**：`/opsx-archive vX.Y.Z` 表示将活跃 changes 收束进目标版本基线，**不**保留 `openspec/changes/archive/` 日期目录。权威入口：`python hack/archive-openspec-to-version.py vX.Y.Z`。收版后必须将本文件中的基线路径更新为目标版本。

### 单服务部署约定（强制）

- 进程模型为**单服务** `illness-service`，**禁止**按微服务拆分为多业务进程；**禁止**在 Compose 中引入 Redis、RabbitMQ 或 Nginx 容器。
- 部署入口：`manifest/docker/docker-compose.yml`（基线）+ `docker-compose.test.yml` / `docker-compose.prod.yml`（overlay）；敏感信息仅经 `.env*`（模板见 `manifest/docker/.env.example`）。
- 镜像发版：`.github/workflows/docker-acr.yml`；`IMAGE_TAG` 由人工控制（起步可为 `v0.0.1`），与 git tag / 服务器 `.env` 对齐。
- **HTTPS / Nginx**：不在本仓库维护 Nginx 或 TLS 配置；边缘反代由宿主机运维侧处理。Go 服务仅监听明文 HTTP（默认 `:9686`）。
- **本地图片 volume**：不为上传目录声明业务挂卷；上传改 OSS 另开变更。

### 数据库连接约定（强制）

- 库名 **`illness`**；可与 `go_ai_talk` 使用同一 MySQL 实例，但 MUST NOT 共享其业务表。
- 连接串 MUST 经 **`ILLNESS_DB_LINK`**（兼容 `GF_DATABASE_DEFAULT_LINK`）注入；可选 **`MYSQL_TCP_HOST`** 覆盖 DSN 主机。启动入口经 `internal/platform/dbcfg` 写入。
- `manifest/config/config.yaml` MUST NOT 提交真实数据库密码。

### OSS / CDN 约定（强制）

- 与 `go_ai_talk` **同桶同 CDN**（当前约定：`bucket=caidoukeji`，`cdnBaseUrl=https://resorce.cuplay.top`）；本仓 `objectKeyPrefix` MUST 为 **`illness/`**。
- AccessKey MUST 经 **`ILLNESS_OSS_ACCESS_KEY_ID` / `ILLNESS_OSS_ACCESS_KEY_SECRET`**（及可选 `ILLNESS_OSS_ENDPOINT`）注入；yaml 中密钥字段留空。
- **上传 HTTP 接口改造**（服务端直传 OSS、返回 CDN URL）MUST 另开 OpenSpec change；未完成前不得以本约定为由擅自改现有 `/upload` 契约。

### 小游戏 / 语音约定

- 微信小游戏首版 **不依赖** `wantSpeak` / `cancelSpeak` / UDP 语音排麦闭环；既有接口可保留 stub，但 MUST NOT 要求新客户端实现。
- 教学局与阶段引导见 OpenSpec capability `tutorial-room-protocol` / `game-phase-guidance`（收版后以当前基线为准）。

## 重要约束

- 评审检查项必须包含“是否引用并遵循当前基线 `openspec/specs/v0.0.1/spec.md`”检查：涉及行为变更时必须可追溯到对应 Requirement/Scenario（基线为空壳时，追溯到本变更的 delta specs 与本文件约定）。
- 评审检查项必须包含“配置是否去密钥、DB/OSS 是否经 env 注入、Compose 是否仍为单服务且无 Redis/MQ/Nginx”。
- 在没有特别要求的情况下，不用生成关于当前变更需求的 md 文件。当有要求生成文档时，文档必须在 `docs/` 文件夹内生成 md 文件。
