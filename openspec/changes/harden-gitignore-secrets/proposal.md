## Why

仓库即将（或正在）首次纳入版本控制时，`.gitignore` 对密钥目录与调试产物覆盖不足：`manifest/docker/env/` 下的真实 `.env.prod` 等易被误提交；`__debug_bin.exe` 等调试二进制已出现在暂存区。需把「本地密钥仓」与构建/调试垃圾文件明确排除，避免密钥与无关二进制进入 git。

**对照说明**：已读 `openspec/project.md` 与空壳基线 `v0.0.1`（首版空壳，本次为 ADDED）。与部署约定「敏感信息仅经 `.env*`、仅提交 `.env.example`」一致并加固。

## What Changes

- **定死忽略** `manifest/docker/env/` 整目录（本地密钥与环境文件放置处；模板仍用 `manifest/docker/.env.example`）
- 补强通用忽略：`*.exe`、`__debug_bin*`、日志、覆盖率等构建/调试产物
- 保留例外：`!.env.example` / `!manifest/docker/.env.example`（及等价「仅示例可提交」规则）
- 实现阶段：对已 staged 的应忽略文件执行 `git rm --cached`（如 `__debug_bin.exe`），不删除工作区文件

### Non-goals

- 不改 Compose / ACR / 业务代码
- 不强制轮换已暴露密钥（可在 tasks 提醒人工处理）
- 不调整 `.vscode/` 是否入库的团队偏好（维持现状忽略）

## Capabilities

### New Capabilities

- `repo-gitignore-hygiene`：仓库忽略规则——`manifest/docker/env/` 整目录、密钥 env 模板例外、调试/构建产物不入库

### Modified Capabilities

- （无）基线 v0.0.1 为空壳

## Impact

- 影响：根目录 `.gitignore`；可能 `git rm --cached` 暂存区中的 exe/env
- 不影响：运行时行为、OpenSpec 基线内容本身（规格仅约束忽略策略）
