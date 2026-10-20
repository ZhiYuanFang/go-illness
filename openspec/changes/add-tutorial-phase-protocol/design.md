## Context

游戏逻辑集中在 `internal/controller/next_illness.go`（内存 `roomArr`）。客户端将改为小游戏并需要可驱动的教学步骤。基线 `v0.0.1` 为空壳。

## Goals / Non-Goals

**Goals:**

- 可创建教学房并按脚本推进，WS 携带引导字段
- 正式写清规则契约，消除 help 与实现偏差
- 向后兼容：旧客户端忽略新 JSON 字段仍可玩娱乐/天梯

**Non-Goals:**

- Redis 房间、语音实现、排行榜、OSS 上传改造、独立 createTutorialRoom 路由

## Decisions

### D1：教学入口形态（已锁定）

- **选择**：仅扩展 `POST /createRoom` 的 `IsTutorial: true`；**不**新增 `createTutorialRoom`。
- 教学局人数 **固定为 4**（1 玩家 + 3 bot）；服务端忽略/覆盖客户端 `RoomNumber`；强制 `IsLyingType=false`。

### D2：脚本推进与 TutorialGate（已锁定）

- **选择**：显式 `tutorialStep` 整数；步骤表为**包内常量**（step → 允许动作集合 + tip + suggestedAction）。
- **门禁结构**（强制，避免分叉）：

```
HTTP 动作 handler
  → TutorialGate(room, action)   // 仅 isTutorial 生效；非法则直接错误返回
  → 原有业务逻辑（同一 night/day 状态机）
  → 若教学且本步成功：tutorialStep++，刷新 tip/suggestedAction
  → BroadcastRoom
```

- 禁止在各 handler 内散落互不相同的 `if isTutorial && step!=N` 而无统一入口。

### D3：引导字段（已锁定）

推送房间 JSON 增加：

| 字段 | 类型 | 说明 |
|------|------|------|
| `isTutorial` | bool | 是否教学局 |
| `tutorialStep` | int | 当前步骤 |
| `phaseTip` | string | 中文短提示 |
| `suggestedAction` | string | `none\|verify\|sleep\|quarantine\|prepare\|tag\|wait` |

教学局 MUST 带齐上述字段；普通局 MAY 省略或给空。旧端忽略未知字段。首版枚举以上表为准；小游戏若需新动作，另开 change 扩枚举。

### D4：教学完成与跳过（已锁定）

- `user` 表增加 `tutorial_completed`（bool/tinyint）；DAO/entity/UserInfo 同步。
- **胜利**：教学局正常胜利结算时置 true。
- **跳过**：独立 **`POST /skipTutorial`**（鉴权用户）；将标记置 true；若用户在教学房中可顺带解散/退出教学房（实现选简单路径：置标记即可，房间按现有离开逻辑清理）。
- 否定：仅客户端 localStorage。

### D5：感染规则契约（已锁定）

规格写明：感染在夜间流程由服务端施加；**无**开局预感染。若未来要开局感染，另开 change。

### D6：语音（已锁定）

mini-game 非目标；不删除既有字段/stub，不扩展。

### D7：步骤表形态（已锁定）

代码常量即可（非 DB、非配置中心）。apply 时在 `next_illness` 旁或 `internal/.../tutorial_script.go` 放置步骤表；至少覆盖：开局认识身份 → 夜查验 → 日隔离/准备 等可玩通的最短路径（具体步数实现时定，但 MUST 每步有 tip + suggestedAction + 允许动作集）。

## Risks / Trade-offs

- [教学与正常局分叉] → 统一 TutorialGate + 复用状态机
- [user 表迁移] → tasks 附 ALTER SQL；库名 `illness`
- [createRoom 参数] → IsTutorial 时强制人数与非 lying

## Migration Plan

1. 模型 / API / user 字段 + ALTER
2. TutorialGate + 步骤常量 + createRoom(IsTutorial)
3. WS 字段与广播
4. skipTutorial + 胜利写完成标记
5. 与小游戏联调；旧端忽略新字段回归

## Open Questions

- （无）推荐项已全部锁定；步骤表具体文案与步数在 apply 时按最短可通关路径填写，不阻塞本提案。
