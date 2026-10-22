package illness

//房间信息
type RoomInfo struct {
	IsSinglePlay   bool                  `json:isSinglePlay`     //是否为单机
	RoomId         int64                 `json:"roomId"`         //房间ID
	UdpPort        int                   `json:"udpPort"`        //UDP端口号
	PlayerArr      map[int64]*PlayerInfo `json:"playerArr"`      //房间内玩家信息
	PlayerJSONArr  []*PlayerInfo         `json:"playerJSONArr"`  //房间内玩家信息
	AllIdentifyArr []int                 `json:"allIdentifyArr"` //房间内所有的身份信息
	// DayTime              int                   //白天时长 秒
	NightTime             int           `json:"nightTime"`             //夜晚时长 秒
	SpeakTime             int           `json:"speakTime"`             //说话的时长 秒
	VoiceTurnRunArr       []*PlayerInfo `json:"voiceTurnRunArr"`       //排麦的人
	GameStat              int           `json:"gameStat"`              //房间当前状态
	CurrentNightTimeStamp int64         `json:"currentNightTimeStamp"` //当前晚上开始的时间点
	CurrentSpeakTimeStamp int64         `json:"currentSpeakTimeStamp"` //当前说话开始的时间点
	IsSuccess             bool          `json:"isSuccess"`             //当前是赢还是输
	NumNight              int           `json:"numNight"`              //当前是第几晚
	LimitRoomNumber       int           `json:"limitRoomNumber"`       //房间人数
	IsLyingType           bool          `json:"isLyingType"`           //是否说谎模式  默认推理模式
	// 教学局引导字段（旧客户端可忽略未知 JSON）
	IsTutorial      bool   `json:"isTutorial"`      // 是否教学局
	TutorialStep    int    `json:"tutorialStep"`    // 当前教学步骤（从 0 起）
	PhaseTip        string `json:"phaseTip"`        // 面向玩家的中文短提示
	SuggestedAction string `json:"suggestedAction"` // none|verify|sleep|quarantine|prepare|tag|wait
	PhaseEndsAt     int64  `json:"phaseEndsAt"`     // 当前相位预计结束 unix 秒；0 表示无时限
}

/*
GameStat
0 默认不考虑
1 房间开启 --- 单一用户
2 身份分发完成 --- 单一用户
3 入夜（玩家行动） 携带倒计时--- 房间用户 【后续会跟随3状态  玩家发动技能调用http接口请求】
4 白天 携带倒计时--- 房间用户
5 排麦【NowSpeekPlayer 就是现在要开麦说话的人】
-1 房间结束



*/
