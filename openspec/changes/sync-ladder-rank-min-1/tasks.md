## 1. 钳位实现

- [x] 1.1 抽出 `RankMin = 1`（或等价常量），替换魔法数 3
- [x] 1.2 修改 `SyncLadder`：失败下限 1；更新注释与 `api/v1` 说明
- [x] 1.3 修改 `UpDateUserRank` 与 `next_illness` 遗留结算钳位为同一下限

## 2. 联调

- [x] 2.1 与前端 `illness` / `ladder-outcome-rank-gate` 联调：连负至 1 再负仍为 1；胜 +1
- [x] 2.2 确认响应 JSON 的 `rank`/`maxRank` 与 DB 一致
