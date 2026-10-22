package controller

import (
	"errors"
	"time"

	"hello/internal/consts"
	"hello/internal/model/illness"
)

// 教学建议主操作枚举（与 OpenSpec game-phase-guidance 一致）
const (
	SuggestedNone       = "none"
	SuggestedVerify     = "verify"
	SuggestedSleep      = "sleep"
	SuggestedQuarantine = "quarantine"
	SuggestedPrepare    = "prepare"
	SuggestedTag        = "tag"
	SuggestedWait       = "wait"
)

// 教学门禁动作名（与 HTTP 能力对应）
const (
	TutorialActionPrepare    = "prepare"
	TutorialActionVerify     = "verify"
	TutorialActionQuarantine = "quarantine"
	TutorialActionSleep      = "sleep"
	TutorialActionTag        = "tag"
	TutorialActionWait       = "wait"
)

// TutorialStepDef 单步教学脚本：允许动作 + 引导文案。
type TutorialStepDef struct {
	AllowedActions  []string
	PhaseTip        string
	SuggestedAction string
}

// tutorialSteps 最短可通关路径（代码常量，非 DB）。
// 0 准备 → 1 夜查验 → 2 日隔离（正确隔离感染者即胜利）。
var tutorialSteps = []TutorialStepDef{
	{
		AllowedActions:  []string{TutorialActionPrepare},
		PhaseTip:        "教学：点下方「准备」。全员就绪后将发身份并入夜。",
		SuggestedAction: SuggestedPrepare,
	},
	{
		AllowedActions:  []string{TutorialActionVerify},
		PhaseTip:        "教学：点头像查验感染。真医生看真相；天梯下庸医结果会取反。",
		SuggestedAction: SuggestedVerify,
	},
	{
		AllowedActions:  []string{TutorialActionQuarantine},
		PhaseTip:        "教学：点头像隔离嫌疑人。误隔健康者即败；隔对全部感染者则胜。",
		SuggestedAction: SuggestedQuarantine,
	},
}

// GetTutorialStep 返回步骤定义；越界返回 nil（教学脚本结束）。
func GetTutorialStep(step int) *TutorialStepDef {
	if step < 0 || step >= len(tutorialSteps) {
		return nil
	}
	return &tutorialSteps[step]
}

// ApplyTutorialTip 按当前 TutorialStep 刷新房间 tip / suggestedAction。
func ApplyTutorialTip(room *illness.RoomInfo) {
	if room == nil || !room.IsTutorial {
		return
	}
	def := GetTutorialStep(room.TutorialStep)
	if def == nil {
		room.PhaseTip = "教学完成：可查看结算或返回。"
		room.SuggestedAction = SuggestedNone
		return
	}
	room.PhaseTip = def.PhaseTip
	room.SuggestedAction = def.SuggestedAction
}

// ApplyPhaseGuidance 教学优先；否则按 GameStat 填充天梯/娱乐局 tip。
func ApplyPhaseGuidance(room *illness.RoomInfo) {
	if room == nil {
		return
	}
	if room.IsTutorial {
		ApplyTutorialTip(room)
		return
	}
	switch room.GameStat {
	case consts.RoomWait:
		room.PhaseTip = "在等：全员准备后开始。"
		room.SuggestedAction = SuggestedPrepare
	case consts.IdentifySendState:
		room.PhaseTip = "身份已下达，请确认你的角色牌。"
		room.SuggestedAction = SuggestedWait
	case consts.NightState:
		room.PhaseTip = "入夜：医生点头像查验；平民等待天亮。"
		room.SuggestedAction = SuggestedVerify
	case consts.DayState:
		room.PhaseTip = "白天：点头像隔离你认定的感染者（误隔即败）。"
		room.SuggestedAction = SuggestedQuarantine
	case consts.GameOver:
		room.PhaseTip = "本局结束。"
		room.SuggestedAction = SuggestedNone
	default:
		if room.PhaseTip == "" {
			room.PhaseTip = "在等：阶段同步…"
			room.SuggestedAction = SuggestedWait
		}
	}
}

// RefreshPhaseEndsAt 按当前 GameStat 写入相位结束时刻。
func RefreshPhaseEndsAt(room *illness.RoomInfo) {
	if room == nil {
		return
	}
	now := time.Now().Unix()
	switch room.GameStat {
	case consts.NightState:
		nt := room.NightTime
		if nt <= 0 {
			nt = 8
		}
		if room.CurrentNightTimeStamp > 0 {
			room.PhaseEndsAt = room.CurrentNightTimeStamp + int64(nt)
		} else {
			room.PhaseEndsAt = now + int64(nt)
		}
	case consts.DayState:
		st := room.SpeakTime
		if st <= 0 {
			st = 20
		}
		room.PhaseEndsAt = now + int64(st)
	default:
		room.PhaseEndsAt = 0
	}
}

// TutorialGate 教学房统一动作门禁；非教学房直接放行。
func TutorialGate(room *illness.RoomInfo, action string) error {
	if room == nil || !room.IsTutorial {
		return nil
	}
	def := GetTutorialStep(room.TutorialStep)
	if def == nil {
		return errors.New("教学步骤已结束，无法继续该操作")
	}
	for _, a := range def.AllowedActions {
		if a == action {
			return nil
		}
	}
	return errors.New("当前教学步骤不允许该操作，请按提示行动")
}

// AdvanceTutorial 教学步成功后递增步骤并刷新 tip。
func AdvanceTutorial(room *illness.RoomInfo) {
	if room == nil || !room.IsTutorial {
		return
	}
	room.TutorialStep++
	ApplyTutorialTip(room)
}
