## Context

主路径本地对局结算走 `POST /syncLadder`（`cUser.SyncLadder`）：胜 +1 / 负 −1，当前多处硬编码 `rank < 3` 抬到 3。兼容 `UpDateUserRank` 与遗留 `next_illness` 房内结算同样下限 3。前端兄弟仓将等级下限改为 1、席位另用 4，要求本仓同步。

## Goals / Non-Goals

**Goals:**

- 所有服务端天梯 rank 失败钳位统一为 **1**
- 注释/常量与前端 `RANK_MIN=1` 契约一致

**Non-Goals:**

- 不实现席位人数（前端事）
- 不改 Win/TutorialDone 请求体
- 不强制迁移历史用户 rank；不新建 `*_test.go`

## Decisions

### D1：单一常量

```go
// 建议 internal/consts 或 user 包内
const RankMin = 1
```

- `SyncLadder`、`UpDateUserRank`、`next_illness` 结算共用，避免再散落魔法数 3

### D2：SyncLadder 行为

```
rank = max(eUser.Rank, RankMin) // 读入异常低值时抬到 1（可选：仅失败侧钳）
if Win: rank++
else: rank = max(rank-1, RankMin)
maxRank = max(maxRank, rank)
```

- 读入时若历史脏数据 &lt;1，写入前钳到 1 即可

### D3：兼容路径同改

- `UpDateUserRank` 的 `target < 3` → `RankMin`
- `next_illness` 房内战 `rank < 3` → `RankMin`（主路径已是本地引擎+syncLadder，但仍防旧客户端）

### D4：注册缺省

- 不强制改为 1；若 DB/注册仍为 3 或 4，可保留（产品只要求「可掉到 1」）

## Risks / Trade-offs

- [前端未发、后端先发] → 低风险：前端仍按 ≥3/4 展示亦可读 1
- [后端未发、前端先发] → 负到 3 被抬回，联调失败 → 同发

## Migration Plan

1. 改常量与三处钳位 + 注释
2. 与前端 `ladder-outcome-rank-gate` 联调：连负至 1、再负仍 1；胜从 1 升 2
3. 部署 go 后再验小游戏结算

## Open Questions

- 无
