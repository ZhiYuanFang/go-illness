## Context

`gf run` / `go run` 只认进程环境变量；密钥约定放在 `manifest/docker/env/`（不入库）。用户在 Windows PowerShell 下开发，要求默认读 `.env.prod`，并在 README 写清步骤。

## Goals / Non-Goals

**Goals:**

- 一条命令：加载默认 `.env.prod` → `gf run main.go`
- 支持 `-EnvFile` 覆盖路径
- 本机默认补 `ILLNESS_OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com`（未设置时）
- README 用中文写清：复制/准备 env、脚本用法、端口验收、可选教学迁移

**Non-Goals:**

- 不做 bash 版（除非后续需要）
- 不改 Compose、不改业务代码
- 不把真实密钥写入 README 或脚本
- 不强制新建 `.env.local`（默认就是 `.env.prod`）

## Decisions

### D1：默认文件为 `.env.prod`（已锁定）

- **选择**：`$EnvFile` 缺省为 `manifest/docker/env/.env.prod`。
- **理由**：用户明确要求；当前仓库本地已在使用该文件。
- **备选**：`.env.local` 优先——已否决。

### D2：仅 PowerShell 脚本

- **选择**：`hack/run-local.ps1`。
- **理由**：当前开发机为 Windows；保持最小面。

### D3：解析规则

- **选择**：跳过空行与 `#` 开头行；按第一个 `=` 分割 key/value；写入当前进程 `Env:`。
- **理由**：兼容密码中含 `+` 等特殊字符。

### D4：OSS endpoint 默认补全

- **选择**：若进程中尚无 `ILLNESS_OSS_ENDPOINT`，设为 `oss-cn-hangzhou.aliyuncs.com`；提供 `-SkipOssEndpoint` 跳过。
- **理由**：本机非 VPC 时内网 endpoint 不可达；上传测需要。

### D5：README 范围

- **选择**：以本地运行为主（脚本 + env 路径 + `9686` 验收）；可保留极简 Makefile/CLI 提示或删掉过时 k8s 模板句，避免误导。
- **理由**：现 README 仍是 GoFrame 单仓模板，与本仓 Compose/env 约定不符。

## Risks / Trade-offs

- [默认读 prod 可能误连生产库] → Mitigation：README 醒目提示「当前默认 `.env.prod`，确认目标库后再启动」；可用 `-EnvFile` 指到其它文件
- [执行策略禁止脚本] → Mitigation：README 写明 `Set-ExecutionPolicy -Scope Process Bypass` 或 `powershell -File ...`
- [未安装 gf] → Mitigation：README 写明先 `make cli` 或安装 GoFrame CLI

## Migration Plan

1. 合入脚本与 README → 开发者将密钥放在 `manifest/docker/env/.env.prod`（若不存在则从 `.env.example` 复制）
2. `.\hack\run-local.ps1` 启动验证 `/` 与 `/api.json`

## Open Questions

- （无阻塞）
