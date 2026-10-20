package v1

import "github.com/gogf/gf/v2/frame/g"

type IllnessMatchReq struct {
	g.Meta          `path:"/putPlayerIntoMatchePool" method:"post"`
	PlayerId        int64
	LimitRoomNumber int
}
type IllnessMatchRes struct {
}

type IllnessAddMachinReq struct {
	g.Meta `path:"/addMachin" method:"post"`
	RoomId int64
}
type IllnessAddMachinRes struct {
}

type IllnessMachingMatchReq struct {
	g.Meta          `path:"/putMachineIntoMatchePool" method:"post"`
	PlayerId        int64
	LimitRoomNumber int
}
type IllnessMachingMatchRes struct {
}

type IllnessCancelMatchReq struct {
	g.Meta   `path:"/cancelMatch" method:"post"`
	PlayerId int64
}
type IllnessCancelMatchRes struct {
}

type IllnessCreateRoomReq struct {
	g.Meta      `path:"/createRoom" method:"post"`
	Id          int64
	IsLyingType bool
	RoomNumber  int
	IsTutorial  bool // true：教学局（固定 4 人、非 lying）
}
type IllnessCreateRoomRes struct {
}

// IllnessSkipTutorialReq 主动跳过教学并持久化完成标记（小游戏/客户端入口）。
type IllnessSkipTutorialReq struct {
	g.Meta `path:"/skipTutorial" method:"post"`
	Id     int64
}
type IllnessSkipTutorialRes struct {
}

// 语音排麦（wantSpeak/cancelSpeak）为历史 stub：小游戏首版非目标，新客户端不必实现。
type IllnessContinueRoomReq struct {
	g.Meta `path:"/continueRoom" method:"post"`
	Id     int64
	RoomId int64
}
type IllnessContinueRoomRes struct {
}
type IllnessPraticeReq struct {
	g.Meta          `path:"/pratice" method:"post"`
	PlayerIds       []int64
	LimitRoomNumber int
}
type IllnessPraticeRes struct {
}
type IllnessRoomWantSpeakReq struct {
	g.Meta `path:"/wantSpeak" method:"post"`
	Id     int64
	RoomId int64
}
type IllnessRoomWantSpeakRes struct {
}

type IllnessRoomCancelSpeakReq struct {
	g.Meta `path:"/cancelSpeak" method:"post"`
	Id     int64
	RoomId int64
}
type IllnessRoomCancelSpeakRes struct {
}

type IllnessDividePersonReq struct {
	g.Meta   `path:"/dividePerson" method:"post"`
	PlayerId int64
	RoomId   int64
}
type IllnessDividePersonRes struct {
}

type IllnessPersonSleepReq struct {
	g.Meta     `path:"/personSleep" method:"post"`
	PlayerId   int64
	RoomId     int64
	IsAllSleep bool
}
type IllnessPersonSleepRes struct {
}

type IllnessVerInfectedPersonReq struct {
	g.Meta   `path:"/verInfectedPerson" method:"post"`
	Id       int64
	PlayerId int64
	RoomId   int64
}
type IllnessVerInfectedPersonRes struct {
	IsInfected bool `json:"isInfected"`
}

type IllnessReportPersonReq struct {
	g.Meta   `path:"/reportPerson" method:"post"`
	Id       int64
	PlayerId int64
	RoomId   int64
}
type IllnessReportPersonRes struct {
}

type IllnessClearHandingUpPersonReq struct {
	g.Meta `path:"/clearHandingUpPerson" method:"post"`
	Id     int64
	RoomId int64
}
type IllnessClearHandingUpPersonRes struct {
}

type IllnessPlayerPrepareReq struct {
	g.Meta `path:"/preparePlay" method:"post"`
	Id     int64
	RoomId int64
}
type IllnessPlayerPrepareRes struct {
}

type IllnessPlayerStartGameReq struct {
	g.Meta `path:"/startGame" method:"post"`
	RoomId int64
}
type IllnessPlayerStartGameRes struct {
}
type IllnessJoinRoomReq struct {
	g.Meta   `path:"/joinRoom" method:"post"`
	PlayerId int64
	RoomId   int64
}
type IllnessJoinRoomRes struct {
}

type IllnessRemoveMachinReq struct {
	g.Meta   `path:"/removeMachin" method:"post"`
	Id       int64
	PlayerId int64
	RoomId   int64
}
type IllnessRemoveMachinRes struct {
}

type IllnessExitRoomReq struct {
	g.Meta `path:"/exitRoom" method:"post"`
	Id     int64
	RoomId int64
}
type IllnessExitRoomRes struct {
}

type IllnessTagIdentifyReq struct {
	g.Meta       `path:"/tagIdentify" method:"post"`
	PlayerId     int64
	TagIndentify int
	RoomId       int64
}
type IllnessTagIdentifyRes struct {
}

type IllnessTagInfoReq struct {
	g.Meta   `path:"/tagInfo" method:"post"`
	PlayerId int64
	RoomId   int64
	TagInfo  string
}
type IllnessTagInfoRes struct {
}
