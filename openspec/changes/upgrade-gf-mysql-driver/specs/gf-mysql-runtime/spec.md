## ADDED Requirements

### Requirement: GoFrame 与 MySQL 驱动版本对齐

系统 MUST 使用与兄弟仓约定一致的 GoFrame 主线版本（本变更目标为 **v2.6.3**），并 MUST 通过空白导入注册 **`github.com/gogf/gf/contrib/drivers/mysql/v2`**（版本与 gf 主模块一致），使 ORM 能解析 `mysql:` 前缀的 Link 并找到名为 `mysql` 的驱动。

#### Scenario: 驱动已注册

- **WHEN** 进程启动且已加载 MySQL contrib 驱动空白导入
- **THEN** GoFrame 内部驱动表中存在 `mysql` 驱动，不再因「未导入驱动」导致查库失败

#### Scenario: Link 可解析 Type

- **WHEN** 启动时经 `ILLNESS_DB_LINK`（或兼容 env）调用 `dbcfg.ApplyGroupFromEnv` 写入形如 `mysql:...@tcp(...)/illness` 的 Link
- **THEN** 首次业务访问 `database.default` 时配置节点 Type 为 `mysql`（非空串），且可创建 DB 连接（网络/账号错误除外）

### Requirement: 登录路径可到达数据库层

在有效 DB 环境变量与可达 MySQL 前提下，`POST /login`（或等价登录入口）在完成微信侧换票后 MUST 能执行用户表查询/插入，MUST NOT 再返回 `cannot find database driver for specified database type ""`。

#### Scenario: 不再出现空 Type 驱动错误

- **WHEN** 本地已用 `hack/run-local.ps1`（或等价）注入合法 `ILLNESS_DB_LINK` 并启动服务后发起登录
- **THEN** 响应与日志中 MUST NOT 出现 database type 为空字符串的驱动查找错误
