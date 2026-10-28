## ADDED Requirements

### Requirement: syncLadder 失败下限为 1

`POST /syncLadder` 在处理失败结算时，MUST 将用户 `rank` 钳至不低于 **1**；胜利时 MUST 将 `rank` 加 1。服务端 MUST NOT 盲信客户端提交的最终 rank 绝对值（保持既有由 Win 推导 ±1 的模式）。同步后 MUST 按需更新 `maxRank`。

#### Scenario: 从高于下限失败降一级

- **WHEN** 用户当前 rank 为 5 且请求 `Win=false` 的 syncLadder 成功
- **THEN** 响应与持久化 rank MUST 为 4

#### Scenario: 已在下限再失败

- **WHEN** 用户当前 rank 为 1 且请求 `Win=false` 的 syncLadder 成功
- **THEN** 响应与持久化 rank MUST 仍为 1

#### Scenario: 胜利升一级

- **WHEN** 用户当前 rank 为 1 且请求 `Win=true` 的 syncLadder 成功
- **THEN** 响应与持久化 rank MUST 为 2，且 maxRank MUST 不小于 2

### Requirement: 兼容写 rank 路径同等下限

仍保留的 `UpDateUserRank`（及遗留房内战结算若仍写 rank）在钳位目标 rank 时，MUST 使用与 syncLadder 相同的下限 **1**，MUST NOT 继续使用下限 3。

#### Scenario: 兼容接口不低于 1

- **WHEN** 兼容更新路径计算出的目标 rank 低于 1
- **THEN** 持久化 rank MUST 被钳为 1
