# 仓库级 AI 执行约定

## OpenSpec

- 全局约定与当前基线指针见 **`openspec/project.md`**。
- 发起 `/opsx-propose`、`/opsx-apply` 前 MUST 先阅读 `openspec/project.md` 及其声明的当前基线 `openspec/specs/vX.Y.Z/spec.md`。
- 收版使用 `/opsx-archive vX.Y.Z`（权威脚本：`hack/archive-openspec-to-version.py`），不是上游「单 change 日期归档」。

## 文档语言

- OpenSpec 相关文档（proposal/design/specs/tasks）默认使用中文说明性文本。
- 环境变量、路径、接口、协议关键字、代码符号可保留英文原文。
