## Why

本仓 OpenSpec 仍沿用上游默认语义（单 change 归档到 `changes/archive/`），缺少全局约束文件与按版本收束的基线规格。兄弟仓 `go_ai_talk` 已验证「`project.md` + 版本基线 + 收版脚本」能强制每次建需求对照历史规格、评估对全局代码的影响。现在需要对齐该工程架构，并从空壳 `v0.0.1` 起步，避免后续 change 脱离全局上下文。

## What Changes

- 新增 `openspec/project.md` 作为全局描述与强制约定（文档语言、OpenSpec 基线参考、收版语义）；基线指针指向 `openspec/specs/v0.0.1/spec.md`
- 引导空壳首版基线 `openspec/specs/v0.0.1/spec.md`（无 capability 正文，仅版本头与空目录说明）
- 从 `go_ai_talk` 移植收版工具链：`hack/archive-openspec-to-version.py`、`split-openspec-baseline.py`、`merge-openspec-specs.py`，并适配空壳基线（无 capability 时可 split/merge）
- **BREAKING（工作流）**：将 `/opsx-archive` 与 `openspec-archive-change` 从「归档单个 change 到日期目录」改为「按版本号收版」；收版后不保留 `openspec/changes/archive/` 日期目录
- 对齐 Cursor 侧 OpenSpec 技能/命令：新增或覆写 `openspec` 总览 skill、`opsx-archive` command/skill；在 propose/apply 流程中强制先读 `project.md` 与当前版本基线
- 可选新增根目录 `AGENTS.md` 中与 OpenSpec 相关的短指针（指向 `openspec/project.md`），不照搬兄弟仓业务域约定（微服务/Redis/Hub 等）

## Capabilities

### New Capabilities

- `openspec-versioned-workflow`: 版本化 OpenSpec 工程约定——`project.md` 全局描述、`specs/vX.Y.Z` 合并基线、按版本收版、新建需求必须对照当前基线评估全局影响

### Modified Capabilities

- （无）当前 `openspec/specs/` 无已有 capability

## Impact

- 影响范围：`openspec/**`、`hack/*openspec*`、`.cursor/skills/**`、`.cursor/commands/opsx-archive.md`（及必要时 propose/apply 中的基线阅读步骤）、可选根 `AGENTS.md`
- 不影响业务 API、数据库或运行时行为
- 后续所有 `/opsx-propose` / `/opsx-apply` / `/opsx-archive` 行为以本约定为准；首次收版目标版本从 `v0.0.1` 起算，下一版为 `v0.0.2` 等
