## ADDED Requirements

### Requirement: 单进程 Compose 部署

仓库 MUST 提供面向单服务 `illness-service` 的 Docker Compose 基线与 test/prod overlay。栈内 MUST NOT 包含 Redis、RabbitMQ 或 Nginx 服务定义；MUST NOT 为本地图片目录声明业务 volume。

#### Scenario: 本地可用基线 compose

- **WHEN** 开发者使用基线 compose 并完成镜像构建
- **THEN** 可启动唯一业务容器，进程监听配置的 HTTP 端口，且栈中无 redis/mq/nginx 容器

#### Scenario: test/prod 使用仓库镜像

- **WHEN** 使用 test 或 prod overlay 并设置 `REGISTRY` 与 `IMAGE_TAG`
- **THEN** 业务容器从对应镜像仓库拉取镜像运行，而非以服务器现场 `--build` 作为常态发布方式

### Requirement: 环境变量管理敏感信息

仓库 MUST 提供 `manifest/docker/.env.example`（或等价路径）列出部署所需变量占位（含镜像、`ILLNESS_DB_LINK`、OSS 密钥等）。真实 `.env` / `.env.*` 含密钥文件 MUST NOT 提交入库。

#### Scenario: 示例环境文件可复制

- **WHEN** 运维复制 `.env.example` 为本地/服务器 `.env` 并填入真实值
- **THEN** Compose 可通过环境变量向容器注入配置，且示例文件本身不含真实密钥

### Requirement: ACR 自动构建推送

仓库 MUST 提供 GitHub Actions 工作流，在推送符合约定的 `v*` tag（及可选手动触发）时构建并推送 `illness-service` 镜像至阿里云 ACR。镜像主 tag MUST 与部署侧 `IMAGE_TAG` 对齐；版本号由人工控制，起步可为 `v0.0.1`。

#### Scenario: 打 tag 触发构建

- **WHEN** 维护者推送 tag `v0.0.1`（或其它人工选定的 `v*` 版本）
- **THEN** 工作流构建并推送对应 tag 的 illness 镜像到配置的 ACR 命名空间

#### Scenario: 单服务 matrix

- **WHEN** 工作流执行构建矩阵
- **THEN** 仅包含本仓单一业务镜像，不要求多微服务矩阵

### Requirement: 现代化容器镜像构建

Dockerfile MUST 在 CI/本地可完成多阶段（或等价）构建，MUST NOT 依赖提交预编译的宿主机 `main` 二进制作为唯一构建输入。

#### Scenario: CI 无预置二进制即可构建

- **WHEN** 在干净环境仅有源码与 Dockerfile 时执行镜像构建
- **THEN** 可成功产出可运行的 illness 服务镜像
