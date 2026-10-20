package controller

import (
	"errors"

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
		PhaseTip:        "欢迎来到教学局：请点击「准备」，全员就绪后将自动发牌并入夜。",
		SuggestedAction: SuggestedPrepare,
	},
	{
		AllowedActions:  []string{TutorialActionVerify},
		PhaseTip:        "入夜后请查验一名玩家。真医生看到真实感染状态；天梯下庸医结果会取反。",
		SuggestedAction: SuggestedVerify,
	},
	{
		AllowedActions:  []string{TutorialActionQuarantine},
		PhaseTip:        "白天请隔离你认为已感染的人。隔离健康玩家会失败；隔离全部感染者则胜利。",
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
		room.PhaseTip = "教学步骤已完成，可结算或退出。"
		room.SuggestedAction = SuggestedNone
		return
	}
	room.PhaseTip = def.PhaseTip
	room.SuggestedAction = def.SuggestedAction
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
