## Why

微信小游戏客户端（`wx-game/illness`）将复活「病变」并对流程做游戏化改造：需要**服务端权威教学剧本**与**阶段引导事件**，否则客户端只能本地假引导，易与房间状态机打架。现有 `isNew` 未落地、规则文案与实现不一致、语音相关 API 多为 stub。需在不拆微服务、不引入 Redis/MQ 的前提下扩展房间协议，供小游戏消费。

**对照说明**：已读 `openspec/project.md` 与空壳基线 `openspec/specs/v0.0.1/spec.md`（首版空壳，本次为 ADDED）。并存活跃工程 change：`align-openspec-versioned-baseline`、`revive-illness-deploy-baseline`（部署/OSS，不冲突；上传/OSS 非本变更范围）。

## What Changes

- 扩展 `POST /createRoom`：`IsTutorial=true` 创建教学局；**固定 4 人**（1 真人 + 3 bot）、强制 `IsLyingType=false`，忽略/覆盖客户端 `RoomNumber`
- 教学房维护 `tutorialStep`；步骤表为**代码内常量**；经统一 **TutorialGate** 做动作门禁后再走现有夜/日状态机；成功后递增步骤并 WS 广播
- WS 房间快照增加：`isTutorial`、`tutorialStep`、`phaseTip`、`suggestedAction`（枚举见 design）；旧客户端忽略未知字段
- 教学胜利自动置 `tutorial_completed`；另提供 **`POST /skipTutorial`** 主动跳过并持久化同一标记；`UserInfo` 可读该标记
- **规则对齐**：感染仅在夜间流程由服务端施加（无开局预感染）；天梯 rank 胜 +1 / 负 −1 / 下限 3；废弃「失败重置到 3」作为服务端事实
- **首版语音非目标**：保留 `wantSpeak` 等字段/stub，不要求小游戏实现

### Non-goals

- 不引入 Redis/MQ，不改单进程内存房间的跨实例扩展
- 不改造 OSS 上传
- 不实现完整排行榜、不实现客户端
- 不新增独立 `createTutorialRoom` 路由（统一走 `createRoom` + `IsTutorial`）

## Capabilities

### New Capabilities

- `tutorial-room-protocol`：教学局创建、TutorialGate 步骤门禁、跳过/完成标记与 WS tip 字段
- `game-phase-guidance`：对局阶段引导载荷（tip / suggestedAction）约定
- `illness-rules-contract`：身份、感染、隔离、天梯 rank、lying 模式的正式规则契约

### Modified Capabilities

- （无）基线 v0.0.1 为空壳；本次均为 ADDED

## Impact

- 影响：`internal/controller/next_illness.go`（建议抽出 TutorialGate 辅助，避免门禁散落）、`internal/model/illness/*`、`api/v1/illness.go` / user 相关、WS 推送 JSON、`user` 表字段
- 客户端契约消费者：`wx-game/illness`（未知字段忽略以兼容旧端）
- 依赖：单服务内存房间；JWT/`Token` 鉴权不变
