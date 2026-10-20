## ADDED Requirements

### Requirement: 数据库经环境变量注入

系统 MUST 通过环境变量 `ILLNESS_DB_LINK`（可兼容 `GF_DATABASE_DEFAULT_LINK`）为 GoFrame `database.default` 提供连接串。仓库内 `manifest/config` 的 yaml MUST NOT 提交真实数据库密码或生产连接串。

#### Scenario: 容器启动使用 env 连接串

- **WHEN** 进程启动且 `ILLNESS_DB_LINK` 已设置为指向库 `illness` 的合法连接串
- **THEN** 应用使用该连接串访问 MySQL，而不依赖 yaml 中的明文密钥

#### Scenario: 配置文件无真实密钥

- **WHEN** 审查仓库中的默认 `config.yaml`
- **THEN** 其中 database link 为空、占位或注释说明，不含真实密码

### Requirement: 库名与实例边界

业务数据库名 MUST 为 `illness`。该库与 `go_ai_talk` 各业务库（如 `ai_voice_*`）逻辑隔离；允许部署在同一 MySQL 实例，MUST NOT 要求共享对方表结构。

#### Scenario: 连接目标库为 illness

- **WHEN** 运维按 `.env.example` 填写 `ILLNESS_DB_LINK`
- **THEN** 连接串中的数据库名为 `illness`（同实例其它库名不在本服务默认连接范围内）

### Requirement: HTTP 监听端口

illness 服务默认 HTTP 监听地址 MUST 为 `:9686`（除非部署环境变量显式覆盖且文档说明）。

#### Scenario: 默认端口

- **WHEN** 使用仓库默认配置启动服务且未覆盖 server.address
- **THEN** 服务监听 `9686` 端口
