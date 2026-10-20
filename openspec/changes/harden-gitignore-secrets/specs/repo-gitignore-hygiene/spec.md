## ADDED Requirements

### Requirement: 忽略本地 Docker 环境目录

仓库根 `.gitignore` MUST 包含对 **`manifest/docker/env/`** 整目录的忽略规则。该目录下任意文件（含 `.env.prod`、`.env.test`、误放的示例名）MUST NOT 被纳入版本库。部署用密钥模板 MUST 仅以 `manifest/docker/.env.example`（或仓库明确允许的 example 路径）提交。

#### Scenario: env 目录不被跟踪

- **WHEN** 运维在 `manifest/docker/env/` 下创建含真实密钥的 `.env.prod`
- **THEN** `git status` / `git check-ignore` 显示该路径被忽略，默认 `git add` 不会将其加入索引

### Requirement: 忽略构建与调试产物

`.gitignore` MUST 忽略常见 Go/本地调试产物，至少包括：`*.exe`、名称匹配 `__debug_bin*` 的文件。已误加入索引的此类文件在落地本变更时 MUST 从索引移除（`git rm --cached`），工作区文件可保留。

#### Scenario: 调试二进制不被提交

- **WHEN** 本地出现 `internal/cmd/__debug_bin.exe` 或根目录 `main.exe`
- **THEN** 这些路径被 ignore；若曾 staged，落地后不再出现在待提交变更中（除非强制 `-f`）

### Requirement: 示例环境文件可提交

根目录或 `manifest/docker/.env.example` MUST 仍可通过例外规则提交，且 MUST NOT 含真实密钥。

#### Scenario: example 仍可见

- **WHEN** 审查仓库中的 `manifest/docker/.env.example`
- **THEN** 该文件可被版本控制跟踪，并作为复制到 `manifest/docker/env/` 的模板
