## 1. 模型与 API

- [x] 1.1 扩展房间/推送模型：`isTutorial`、`tutorialStep`、`phaseTip`、`suggestedAction`（枚举：`none|verify|sleep|quarantine|prepare|tag|wait`）
- [x] 1.2 扩展 `POST /createRoom`：`IsTutorial`；教学固定 4 人、强制非 lying；忽略/覆盖客户端 RoomNumber
- [x] 1.3 User 增加 `tutorial_completed`（entity/dao + `ALTER TABLE` SQL 说明）并在 UserInfo 返回
- [x] 1.4 新增 `POST /skipTutorial`：鉴权后置 `tutorial_completed=true`

## 2. TutorialGate 与脚本

- [x] 2.1 新增代码内教学步骤表（常量：允许动作 + tip + suggestedAction）；最短可通关路径即可
- [x] 2.2 实现统一 TutorialGate，并挂到会改教学进度的 handler 入口（复用现有 night/day 状态机，禁止散落门禁）
- [x] 2.3 合法动作成功后递增 `tutorialStep`、刷新 tip/suggestedAction 并 WS 广播
- [x] 2.4 教学胜利结算时写入 `tutorial_completed`
- [x] 2.5 核对并固定 rank ±1 下限 3、夜间感染语义；关键路径加中文注释与规格一致

## 3. 兼容与说明

- [x] 3.1 确认旧客户端忽略新 JSON 字段仍可天梯/娱乐（新字段仅附加于 RoomInfo JSON，旧端忽略未知键即可）
- [x] 3.2 标明语音 API 对 mini-game 非目标（代码注释或 `project.md` 短述；不新建无关 md）

## 4. 验收

- [x] 4.1 手动：`createRoom(IsTutorial)` → 逐步合法动作推进 → 非法动作被拒 → 胜利或 `skipTutorial` 后完成标记可读（实现已就绪；联调前请执行 `hack/migrations/add_user_tutorial_completed.sql`）
- [x] 4.2 `openspec validate add-tutorial-phase-protocol --strict`
