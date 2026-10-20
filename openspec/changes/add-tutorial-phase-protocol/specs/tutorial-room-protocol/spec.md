## ADDED Requirements

### Requirement: 创建教学房间

服务端 MUST 通过扩展 `POST /createRoom` 的 `IsTutorial=true` 创建教学局（MUST NOT 依赖额外专用创建路由作为唯一入口）。教学局 MUST 为服务端权威：**固定人数 4**（1 真人 + 机器人）、`IsLyingType=false`（客户端传入的 lying/人数 MUST 被覆盖或拒绝不一致请求），并按脚本步骤推进。`IsTutorial` 缺失或为 false 时行为 MUST 与现网非教学创建一致。

#### Scenario: 带 IsTutorial 创建

- **WHEN** 已鉴权用户以 `IsTutorial=true` 请求 `createRoom`
- **THEN** 服务端创建 4 人教学房、填充机器人，并通过 WebSocket 向该用户推送含 `isTutorial=true` 的房间快照

### Requirement: TutorialGate 步骤门禁

教学房 MUST 维护 `tutorialStep`；每步允许的动作集合与 tip 来自**代码内步骤表**。所有会改变教学进度的玩家动作 MUST 先经统一 TutorialGate 校验：非法动作 MUST 失败、MUST NOT 推进步骤、MUST NOT 破坏房间一致性。合法动作在业务成功后 MUST 递增 `tutorialStep` 并推送更新快照（含新 tip / suggestedAction）。

#### Scenario: 步骤外行动被拒绝

- **WHEN** 教学局当前步骤仅允许查验，用户却请求隔离
- **THEN** 服务端拒绝该请求并保持 `tutorialStep` 不变

#### Scenario: 合法动作推进步骤

- **WHEN** 教学局当前步骤允许查验且用户成功完成查验
- **THEN** `tutorialStep` 递增且随后推送的快照含更新后的 `phaseTip` 与 `suggestedAction`

### Requirement: 教学完成标记

服务端 MUST 在用户表持久化教学完成标记（如 `tutorial_completed`）。教学局**胜利结算**后 MUST 将该用户标记为已完成。后续 `UserInfo`（或等价）MUST 能让客户端读取该标记。

#### Scenario: 教学胜利后

- **WHEN** 教学局以胜利结束
- **THEN** 该用户的教学完成标记为已完成

### Requirement: 跳过教学接口

服务端 MUST 提供 `POST /skipTutorial`（需鉴权）。调用成功后 MUST 将当前用户教学完成标记置为已完成，且 MUST NOT 要求用户先处于教学房内才能跳过。

#### Scenario: 主动跳过

- **WHEN** 已鉴权用户调用 `skipTutorial`
- **THEN** 其教学完成标记为已完成，且后续 UserInfo 可读到已完成
