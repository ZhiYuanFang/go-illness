## Context

`go-illness` 已具备上游 OpenSpec 脚手架（`openspec/config.yaml`、空 `specs/`、`changes/archive/`、Cursor opsx skills/commands），但语义仍是「单 change → 日期归档」。兄弟仓 `go_ai_talk` 将权威语义改为：活跃 changes 收进 `openspec/specs/vX.Y.Z/spec.md`，以 `openspec/project.md` 为全局宪法并固定当前基线指针。

本仓业务仍为 GoFrame 单仓模板，尚无 illness 域能力规格可沉淀；用户明确要求首版 **`v0.0.1` 空壳基线**，其余 **OpenSpec 工程约束对齐** `go_ai_talk`（不照搬其微服务/Redis/Hub 等业务约定）。

## Goals / Non-Goals

**Goals:**

- 建立 `project.md` + `specs/v0.0.1/spec.md`（空壳）+ 收版脚本 + Cursor 工作流覆写
- 使后续 `/opsx-propose` 必须对照当前基线与 `project.md`，再写增量
- 使 `/opsx-archive vX.Y.Z` 成为收版权威入口；不保留日期 archive 目录

**Non-Goals:**

- 不把当前业务代码行为扫进首版基线
- 不移植 `go_ai_talk` 的领域约定（gateway/Redis/Hub/usage 统计等）
- 不改变业务代码、API、数据库、部署清单
- 不在本变更内执行「收版本 change 自身」以外的版本跳跃

## Decisions

### D1：首版基线形态 — 手工空壳，而非跑收版脚本生成

- **选择**：直接写入 `openspec/specs/v0.0.1/spec.md`，含版本标题、说明引用、空「目录」与「本版无 capability」声明；不含 `## <capability>` 章节。
- **理由**：收版脚本依赖「源基线 + 活跃 changes」；尚无源版本时，bootstrap 必须特殊路径。空壳与后续 merge 产物格式对齐（同头同目录段），便于下一版从 `v0.0.1` 起算。
- **备选**：放一个占位 capability `openspec-bootstrap` —— 拒绝，以免污染「空壳」语义与目录。

### D2：收版脚本从 `go_ai_talk` 移植，并适配空壳

- **选择**：复制三脚本；对空基线：
  - `split-openspec-baseline.py`：解析 0 个 capability 时 **成功退出**（打印 0），不报错；
  - `merge-openspec-specs.py`：无 capability 文件时仍写出仅含头+空目录（或「本版无 capability」）的 `vX.Y.Z/spec.md`；
  - `archive-openspec-to-version.py`：目标为全新版本且无更早源时，若用户指定仅有 `v0.0.1` 空壳，允许以该空壳为源；若无任何版本目录，则要求先存在源基线（bootstrap 不经此脚本）。
- **理由**：与兄弟仓同一入口，仅补空壳缺口，避免本仓卡在「第一次收版」。
- **备选**：空壳期禁止跑脚本 —— 拒绝，下一版收版会立刻撞墙。

### D3：`project.md` 只写 OpenSpec/工程通用约定

- **选择**：包含：项目简述、技术栈（Go/GoFrame）、文档中文强制、代码注释与测试文件约定（与兄弟仓 OpenSpec 段同精神但可按本仓裁剪）、**OpenSpec 基线参考约定（强制）** 指向 `openspec/specs/v0.0.1/spec.md`、收版命令与脚本路径。不写入兄弟仓业务域长文。
- **理由**：用户要求对齐 openspec 约束，不是对齐产品架构。
- **备选**：整文件复制 `go_ai_talk/openspec/project.md` —— 拒绝，会产生错误强制项。

### D4：工作流覆写范围

- **必须**：`.cursor/commands/opsx-archive.md`、`.cursor/skills/openspec-archive-change/SKILL.md` → 按版本收版语义（对齐兄弟仓）。
- **必须**：新增 `.cursor/skills/openspec/SKILL.md` 总览（含收版说明与读基线步骤）。
- **建议**：在 `openspec-propose` / `openspec-apply` skill 顶部增加「先读 `project.md` + 当前基线」硬步骤（可短，不重写整份 skill）。
- **可选**：根 `AGENTS.md` 仅 OpenSpec 指针段落；若创建，保持短小。
- **github/trae 提示**：本仓若无对应树，不强制新建；以 Cursor 侧为准。

### D5：本 change 自身的规格归宿

- **选择**：delta 落在 `openspec/changes/align-openspec-versioned-baseline/specs/openspec-versioned-workflow/spec.md`。实现完成后由用户决定何时 `/opsx-archive v0.0.2`（或用户指定下一版）把本 capability 收进下一版基线；**不**在 apply 时把本 change 提前 merge 进 `v0.0.1`（保持 `v0.0.1` 空壳承诺）。
- **理由**：空壳首版与「本变更引入的工作流能力」分版更清晰。
- **备选**：apply 时直接把 `openspec-versioned-workflow` 写进 `v0.0.1` —— 与「空壳」冲突，拒绝。

### D6：空壳与 OpenSpec CLI 校验

- **选择**：`v0.0.1/spec.md` 为合并视图文件；日常增量仍经 `openspec/changes/<id>/specs/`。`openspec validate` 以 change 为主；合并基线不要求 CLI 把版本目录当 capability。
- **理由**：与兄弟仓一致（版本目录名匹配 `v\d+\.\d+\.\d+`，merge/split 会排除之）。

## Risks / Trade-offs

- [空壳基线无可对照的业务 Requirement] → Mitigation：`project.md` 写明「基线为空时仍须读 project.md；有 capability 后对照最新 `v*`」；后续第一个业务 change 起开始填充。
- [split/merge 空壳行为与兄弟仓分叉] → Mitigation：仅增加「零 capability 成功」分支；有内容后路径与兄弟仓一致；注释标明差异。
- [上游 `openspec archive` 与本仓「收版」易混淆] → Mitigation：command/skill 文案强调参数是 `vX.Y.Z`；禁止 `mv` 到日期 archive。
- [本 change 归档进哪一版] → Mitigation：design 固定 `v0.0.1` 保持空壳；收本 change 用 `v0.0.2+`，tasks 中写清提醒。

## Migration Plan

1. Apply：落地 `project.md`、空壳 `v0.0.1`、三脚本、Cursor 覆写。
2. 验证：`openspec list`；目视基线空壳；dry 思路确认下一版收版入口存在。
3. 不回滚业务代码；若回滚本变更，删除新增文件并恢复上游 archive skill/command 即可。

## Open Questions

- 无阻塞项。若后续希望把 `openspec-versioned-workflow` 直接并入 `v0.0.1`（放弃严格空壳），需用户另行确认后改 design。
