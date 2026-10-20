package illness

import "net"

type PlayerInfo struct {
	UDPZONE             string
	UDPIP               net.IP
	UDPPort             int             //语音的UDP地址端口，玩家连通之后会发送固定格式，记录到本地缓存中，然后服务端收到语音格式 将派发给当前房间的所有玩家
	Id                  int64           `json:"id"` //玩家的ID
	Lv                  float32         `json:"lv"`
	NickName            string          `json:"nickName"`
	Avatar              string          `json:"avatar"`
	StartWaitTime       int64           `json:"startWaitTime"`       //玩家开始匹配等待的时间戳
	Identity            int             `json:"identity"`            //玩家身份 1庸医 2医生 3平民
	GameStat            int             `json:"gameStat"`            //游戏进度控制
	IsQuarantine        bool            `json:"isQuarantine"`        //是否被隔离
	InfectedSeri        int             `json:"infectedSeri"`        //被感染的顺序
	IsInfected          bool            `json:"isInfected"`          //是否已经被感染
	IsSleep             bool            `json:"isSleep"`             //是否睡眠
	TonightVerPlayerId  int64           `json:tonightVerPlayerId`    //医生当夜查询的对象
	BeReportedNum       int             `json:"beReportedNum"`       //被举报的数量
	ReportPlayerIdArr   map[int64]int64 `json:"reportPlayerIdArr"`   //举报此人的玩家ID列表
	HasPrepared         bool            `json:"hasPrepared"`         //是否准备好了
	IsRoomManager       bool            `json:"isRoomManager"`       //是否为房主 目前用不着这个
	CheckContent        []string        `json:"checkContent"`        //查验的结果记录
	WantLimitRoomNumber int             `json:"wantLimitRoomNumber"` //想要匹配的房间人数
	TagIdentify         int             `json:"tagIdentify"`         //标记的身份
	TagInfo             string          `json:"tagInfo"`             //标记文本
}

/*
GameStat
0 默认不考虑
1 房间开启 --- 单一用户
2 身份分发完成 --- 单一用户
3 入夜（玩家行动） 携带倒计时--- 房间用户 【后续会跟随3状态  玩家发动技能调用http接口请求】
4 白天 携带倒计时--- 房间用户
5 排麦？
-1 房间结束



*/
