## Why

本地用 `gf run` 时进程不会自动加载 `.env`，每次手工灌环境变量易错且重复。需要一条可重复的本地启动入口，并把「如何跑本地」写进仓库 `README.MD`，替代仍停留在 GoFrame 模板的说明。

**对照说明**：已读 `openspec/project.md` 与空壳基线 `v0.0.1`（首版空壳，本次为 ADDED）。敏感文件仍仅放 `manifest/docker/env/`（gitignore）；模板仍为 `manifest/docker/.env.example`。用户锁定：**脚本默认加载 `.env.prod`**。

## What Changes

- 新增 `hack/run-local.ps1`：从 `manifest/docker/env/` 加载 env 后执行 `gf run main.go`
- **默认 env 文件**：`manifest/docker/env/.env.prod`（可用参数覆盖）
- 本机开发时若未设置 `ILLNESS_OSS_ENDPOINT`，脚本默认补公网 endpoint（可关闭）
- 重写/补充根目录 `README.MD`：说明本地如何准备 env、如何用脚本启动、如何验收（`:9686`）

## Capabilities

### New Capabilities

- `local-dev-run`：Windows 本地开发启动脚本与 README 本地运行说明

### Modified Capabilities

- （无）基线为空壳；不修改既有未收版 capability 的行为要求

## Impact

- 新增：`hack/run-local.ps1`
- 文档：`README.MD`（替换或大幅更新模板内容）
- 不改业务 Go 代码、Compose、上传契约
- 依赖：本机已安装 `gf` CLI 与 Go；真实密钥由开发者自备于 `env/.env.prod`
