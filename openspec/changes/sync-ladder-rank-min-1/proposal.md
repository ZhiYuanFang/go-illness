## Why

兄弟仓微信小游戏前端（`illness` / change `ladder-outcome-rank-gate`）将天梯**等级账本下限**改为 1，并与**开局席位下限 4**分离。本仓 `SyncLadder` / 兼容 `UpDateUserRank` / 遗留房内战结算仍钳位 **下限 3**，联调会被服务端抬回，前端无法表现 1–2 级。需服务端同步放宽到 1。

## What Changes

- **BREAKING**：天梯失败后 rank 下限由 **3 → 1**（胜 +1、负 −1 语义不变）
- `POST /syncLadder`：读写前后钳位与注释改为 `RANK_MIN=1`
- 兼容路径 `UpDateUserRank` 与遗留 `next_illness` 房内 rank 结算：同等下限 1（避免旧路径仍抬到 3）
- API 注释与常量（若抽出）对齐
- **不做什么**：不负责前端席位 `max(4, rank)`；不强制新用户默认 rank=1（可保留现有注册缺省）；不新增接口字段；本变更不要求新建测试文件（遵循 project 约定）

## Capabilities

### New Capabilities

- `sync-ladder-rank-min`：结算同步与兼容路径的 rank 下限 1 契约

### Modified Capabilities

- （基线 `openspec/specs/v0.0.1/spec.md` 为空壳，无既有 capability 可改；本次为 ADDED）

## Impact

- 已对照 `openspec/project.md` 与基线 **v0.0.1（空壳）**
- 代码：`internal/controller/user.go`（`SyncLadder`、`UpDateUserRank`）、`api/v1/user.go` 注释、`internal/controller/next_illness.go` 遗留结算钳位
- 联调：与前端 `illness` 的 `ladder-outcome-rank-gate` 同发；验收 rank=1 再负仍为 1
