---
name: openspec
description: >-
  指导基于规格的变更提案、增量规格、校验与按版本收版，配合 OpenSpec CLI 与本仓库目录结构。
  在用户提到 OpenSpec、变更提案、规格增量、openspec validate、规划功能、破坏性变更、
  架构调整，或要求创建/应用规格与变更时使用。
---

# OpenSpec（本仓库）

## 权威文档

编写或修改 OpenSpec 产物前，请先阅读：

1. **[openspec/project.md](../../../openspec/project.md)** — 本仓库全局约定与当前基线指针。
2. 当前基线：`openspec/project.md` 中声明的 `openspec/specs/vX.Y.Z/spec.md`（首版空壳为 `v0.0.1`）。
3. 根目录 **[AGENTS.md](../../../AGENTS.md)** — OpenSpec 短指针。

## 何时走 OpenSpec

**创建变更**（提案 + 增量 + `tasks.md`，并用 `--strict` 校验）适用于：新增能力、破坏性 API/架构/安全模式变更、或改变行为的性能优化。

**可跳过提案**的情况：恢复规格所描述行为的缺陷修复、错别字与排版、非破坏性依赖升级、仅配置调整——除非需求模糊，走提案更安全。

## 助手工作流（精简）

### 编写或实现之前（强制）

- **MUST** 先阅读 `openspec/project.md` 与其指向的当前基线 `spec.md`（即便为空壳也须确认）。
- 运行 `openspec list` 了解进行中的变更；检查 `openspec/changes/` 是否已有重叠变更。
- 基线有 capability 时，对照相关章节评估对全局代码的影响。

### 创建变更

1. 选取唯一的动词开头 **change-id**（kebab-case：`add-…`、`update-…`、`remove-…`、`refactor-…`）。
2. 在 `openspec/changes/<change-id>/` 下搭建：`proposal.md`、`tasks.md`、可选 `design.md`，以及 `specs/<capability>/spec.md` 增量。
3. 增量使用 `## ADDED|MODIFIED|REMOVED|RENAMED Requirements`；每条 Requirement 至少包含一个 `#### Scenario:` 块（四个 `#`）。
4. 分享前运行 **`openspec validate <change-id> --strict`** 并修复全部问题。
5. 提案**未经评审通过前不要开始实现**。

### 实现已批准的变更

按顺序阅读 `proposal.md` → `design.md`（若有）→ `tasks.md`；完成任务后，如实勾选 `tasks.md` 中的复选框。实现须对照当前基线 Requirement/Scenario（基线非空时）。

### 收版（/opsx-archive）

将活跃 change 收进目标版本基线（如 `v0.0.2`），**不**保留 `openspec/changes/archive/`：

```bash
python hack/archive-openspec-to-version.py v0.0.2
```

收版后更新 `openspec/project.md` 基线引用，并跑 `openspec validate --strict`。

**注意**：参数是 **`vX.Y.Z` 版本号**，不是单个 change-id；勿按上游「归档单个 change 到 archive/」流程执行。

## CLI 速查

```bash
openspec list                  # 进行中的变更
openspec list --specs          # 规格列表
openspec show <item>           # 查看变更或规格
openspec validate <id> --strict
openspec archive <change-id> [--skip-specs] [--yes]   # 通常由收版脚本调用
openspec update [path]         # 刷新说明文件
```

## 常见坑

- **MODIFIED 需求**：粘贴**完整**更新后的需求块；增量若不全，收版时会丢细节。
- **Scenario**：仅使用 `#### Scenario: 标题`。
- **空壳基线**：`v0.0.1` 无 capability 属预期；首个业务/工程 capability 经 change ADDED，再用 `/opsx-archive v0.0.2` 等收束。
- **检索**：优先用 `openspec` 子命令，并对 `openspec/specs`、`openspec/changes` 使用 `rg`。
