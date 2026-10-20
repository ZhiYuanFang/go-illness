## ADDED Requirements

### Requirement: 以 project.md 为全局 OpenSpec 约束

系统 MUST 在仓库中维护 `openspec/project.md`，作为 OpenSpec 与协作的全局描述文件，至少包含：项目简述、技术栈、文档语言约定、当前基线路径（形如 `openspec/specs/vX.Y.Z/spec.md`）、以及「新建变更前必须对照当前基线」的强制条款。

#### Scenario: 仓库存在全局描述与基线指针

- **WHEN** 贡献者或 AI 发起任意新 OpenSpec 变更
- **THEN** 必须先读取 `openspec/project.md`，并按其声明的当前基线路径打开对应 `spec.md`

#### Scenario: 文档语言为中文

- **WHEN** 生成或修改 OpenSpec 产物（proposal、design、specs、tasks）
- **THEN** 说明性文本 MUST 使用中文；仅环境变量名、路径、接口路径、技术关键字等可保留英文

### Requirement: 版本化合并基线

系统 MUST 将已收束的规格存放于 `openspec/specs/vX.Y.Z/spec.md`（`v` + 三段数字）。版本目录 MUST 与临时 capability 目录区分；merge/split 工具 MUST 将匹配 `^v\d+\.\d+\.\d+$` 的目录视为版本基线而非 capability。

#### Scenario: 首版空壳基线 v0.0.1

- **WHEN** 完成本工作流落地且尚无业务 capability 需并入首版
- **THEN** 仓库 MUST 存在 `openspec/specs/v0.0.1/spec.md`，且该文件为无 capability 章节的空壳（可含版本标题、说明与空目录声明），同时 `project.md` 基线指针 MUST 指向该文件

#### Scenario: 历史版本可查阅

- **WHEN** 曾执行多次按版本收版
- **THEN** 各历史版本目录 `openspec/specs/v*/spec.md` MUST 保留，供新建需求对照全局影响

### Requirement: 按版本收版而非日期归档

`/opsx-archive`（及对应 archive skill）MUST 将参数解释为目标版本号 `vX.Y.Z`，权威执行入口 MUST 为 `hack/archive-openspec-to-version.py`。收版 MUST 把活跃 `openspec/changes/` 下变更合并进目标版本基线，MUST NOT 以「移动到 `openspec/changes/archive/YYYY-MM-DD-<name>` 并保留」作为完成态。收版成功后 MUST 更新 `openspec/project.md` 中的基线指针。

#### Scenario: 带版本号收版

- **WHEN** 用户执行 `/opsx-archive v0.0.2`（或其它合法 `vX.Y.Z`）且存在源基线
- **THEN** 系统将活跃 changes 收束为 `openspec/specs/v0.0.2/spec.md`，清空活跃 change 目录（不保留日期 archive 内容），并更新 `project.md` 基线引用

#### Scenario: 缺少版本号时拒绝按 change-id 归档

- **WHEN** 用户调用 archive 未提供合法 `vX.Y.Z`
- **THEN** 助手 MUST 列出已有版本目录与活跃 changes，并要求用户给出目标版本；MUST NOT 将参数当作单个 change-id 执行日期目录归档

### Requirement: 新建需求必须对照历史基线评估全局影响

AI 在 propose 与 apply 阶段 MUST 对照当前 `project.md` 指向的基线规格；若需求触及已有 capability，proposal MUST 标明复用或变更了哪些基线边界。实现若与基线 Requirement/Scenario 冲突，MUST 先更新变更增量规格再改代码。

#### Scenario: propose 前读取基线

- **WHEN** 执行 `/opsx-propose` 创建新变更
- **THEN** 助手 MUST 先阅读 `openspec/project.md` 与当前基线 `spec.md`（即便基线为空壳，也须确认指针与空壳状态），再生成 proposal/design/tasks/增量 specs

#### Scenario: 基线为空壳时的行为

- **WHEN** 当前基线为无 capability 的空壳（如 v0.0.1）
- **THEN** 新建变更仍 MUST 遵守 `project.md` 中的全局约定，并可以 ADDED 方式引入新 capability；不得因基线为空而跳过阅读 `project.md`

### Requirement: 收版工具链可处理空壳基线

仓库 MUST 提供 `hack/archive-openspec-to-version.py`、`hack/split-openspec-baseline.py`、`hack/merge-openspec-specs.py`。在源基线无 capability 或暂无 capability 目录时，split 与 merge MUST 成功完成（写出或保留合法空壳合并文件），不得因「零 capability」而非零退出（除非路径缺失等真实错误）。

#### Scenario: 从空壳源基线收版到下一版

- **WHEN** 源基线 `v0.0.1` 无 capability，且存在待收版的活跃 change
- **THEN** 运行 `python hack/archive-openspec-to-version.py v0.0.2` MUST 能完成 split → archive → merge，生成含新 capability 的 `openspec/specs/v0.0.2/spec.md`
