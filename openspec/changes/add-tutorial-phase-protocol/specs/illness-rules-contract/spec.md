## ADDED Requirements

### Requirement: 身份与查验

身份取值 MUST 为：1 庸医、2 医生、3 平民。天梯（`IsLyingType=false`）下，非真医生玩家的查验结果 MUST 相对真实感染状态取反；真医生 MUST 看到真实结果。娱乐 lying 模式下，真人庸医查验语义 MUST 按现有实现保持（信息正确、身份自知），并在规格实现注释中写清。

#### Scenario: 天梯庸医查验取反

- **WHEN** 天梯房中庸医查验一名未感染玩家
- **THEN** 返回给该庸医的感染结论为相反值（显示为已感染）

### Requirement: 感染与隔离胜负

感染 MUST 在夜间流程由服务端施加（当前契约：**无**开局预感染）。隔离未感染者 MUST 导致好人失败；全部感染者被隔离且仍存在未感染者 MUST 胜利；全员感染 MUST 失败。

#### Scenario: 误隔离健康玩家

- **WHEN** 白天隔离一名 `isInfected=false` 的玩家
- **THEN** 游戏进入失败结算

### Requirement: 天梯等级

天梯结算对房主 rank：胜利 MUST +1；失败 MUST −1 且 MUST NOT 低于 3。MUST NOT 将「失败直接重置为 3」作为服务端行为，除非未来变更明确修改本条。

#### Scenario: rank 为 3 时失败

- **WHEN** rank=3 的房主天梯失败
- **THEN** rank 仍为 3

### Requirement: 小游戏语音非目标

本变更 MUST 将语音排麦/UDP 标为小游戏首版非目标：不新增语音行为；既有 stub 可保留但 MUST NOT 要求新客户端实现。

#### Scenario: 文档声明

- **WHEN** 实现者阅读本 capability
- **THEN** 可明确得知 mini-game 不依赖 `wantSpeak`/`UDP` 完成闭环
