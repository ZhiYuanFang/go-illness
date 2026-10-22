## ADDED Requirements

### Requirement: 本地启动脚本加载 env 并运行服务

仓库 MUST 提供 Windows PowerShell 脚本 `hack/run-local.ps1`。在未指定覆盖参数时，脚本 MUST 默认从 **`manifest/docker/env/.env.prod`** 加载环境变量（跳过空行与 `#` 注释，按首个 `=` 分割），然后将工作目录置于仓库根并启动 **`gf run main.go`**（或设计约定的等价入口）。脚本 MUST 支持通过参数指定其它 env 文件路径。

#### Scenario: 默认使用 .env.prod

- **WHEN** 开发者在仓库根执行 `.\hack\run-local.ps1` 且存在 `manifest/docker/env/.env.prod`
- **THEN** 脚本将该文件中的键值注入当前进程环境，并启动本地 Go 服务入口

#### Scenario: 覆盖 env 文件

- **WHEN** 开发者传入指向其它已存在 env 文件的参数（如 `-EnvFile`）
- **THEN** 脚本加载该文件而非默认 `.env.prod`

#### Scenario: 默认补全本机 OSS endpoint

- **WHEN** 加载完成后进程环境中尚无 `ILLNESS_OSS_ENDPOINT`，且未选择跳过补全
- **THEN** 脚本将 `ILLNESS_OSS_ENDPOINT` 设为公网 endpoint `oss-cn-hangzhou.aliyuncs.com`

### Requirement: README 说明本地运行

根目录 `README.MD` MUST 用中文说明如何在本地运行本服务，至少包含：环境文件位置（`manifest/docker/env/`，默认 `.env.prod`）、如何从 `.env.example` 准备文件、如何调用 `hack/run-local.ps1`、默认监听端口 `:9686` 与简单验收路径。README MUST NOT 包含真实密钥。

#### Scenario: 新成员按 README 可启动

- **WHEN** 开发者阅读 `README.MD` 的本地运行章节并具备有效的 `env/.env.prod` 与 `gf` CLI
- **THEN** 可按文档步骤启动服务并访问约定端口进行验收
