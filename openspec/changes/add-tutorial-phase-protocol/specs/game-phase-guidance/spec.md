## ADDED Requirements

### Requirement: 房间快照阶段引导字段

向玩家推送的房间 JSON MUST 可包含（教学局 MUST 包含；普通局 MAY 包含）：

- `phaseTip`：面向玩家的中文短提示
- `suggestedAction`：枚举建议主操作（至少支持 `none|verify|sleep|quarantine|prepare|tag|wait`）

旧客户端忽略未知字段时 MUST 仍能解析既有字段并游戏。

#### Scenario: 教学夜提示查验

- **WHEN** 教学局进入需要医生查验的步骤
- **THEN** 推送快照中 `phaseTip` 非空且 `suggestedAction` 为 `verify`

### Requirement: 与 gameStat 一致

`phaseTip` / `suggestedAction` MUST 与当前 `gameStat` 及玩家身份合法动作一致，不得引导玩家执行当前状态禁止的操作。

#### Scenario: 白天不提示夜晚查验

- **WHEN** `gameStat` 为白天
- **THEN** `suggestedAction` MUST NOT 为仅夜晚合法的 `verify`
