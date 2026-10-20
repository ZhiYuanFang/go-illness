## 1. 空壳基线与全局描述

- [x] 1.1 创建 `openspec/specs/v0.0.1/spec.md` 空壳（版本标题、说明、空目录/无 capability 声明，无 capability 章节）
- [x] 1.2 创建 `openspec/project.md`：项目简述、技术栈、中文文档约定、OpenSpec 基线参考约定（指针指向 `openspec/specs/v0.0.1/spec.md`）、收版入口说明；不照搬 `go_ai_talk` 业务域长文
- [x] 1.3 （可选）新增根 `AGENTS.md` 短指针段落，链到 `openspec/project.md`

## 2. 收版工具链

- [x] 2.1 从 `go_ai_talk` 移植 `hack/split-openspec-baseline.py`、`hack/merge-openspec-specs.py`、`hack/archive-openspec-to-version.py`
- [x] 2.2 适配空壳：split 零 capability 成功退出；merge 无 capability 时仍写出合法空壳合并文件
- [x] 2.3 本地验证：对 `v0.0.1` 执行 split（期望 0）与 merge 到临时/同版路径行为符合预期（不破坏空壳约定）

## 3. Cursor 工作流对齐

- [x] 3.1 覆写 `.cursor/commands/opsx-archive.md` 为按版本收版语义（对齐 `go_ai_talk`）
- [x] 3.2 覆写 `.cursor/skills/openspec-archive-change/SKILL.md` 为按版本收版
- [x] 3.3 新增 `.cursor/skills/openspec/SKILL.md` 总览（读 project.md/基线、收版说明）
- [x] 3.4 在 `openspec-propose` / `openspec-apply` skill 中增加「先读 project.md + 当前基线」强制步骤

## 4. 校验与收尾说明

- [x] 4.1 运行 `openspec validate align-openspec-versioned-baseline --strict` 并修复问题
- [x] 4.2 确认 `v0.0.1` 仍为空壳；本 capability 仅存在于本 change 的 delta specs，提醒用户后续用 `/opsx-archive v0.0.2`（或指定下一版）收束本 change
