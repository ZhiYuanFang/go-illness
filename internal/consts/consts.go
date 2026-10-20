package consts

import "github.com/gogf/gf/v2/net/ghttp"

const Authorization = "Authorization"

var PersonWSArr map[int64]*ghttp.WebSocket //WebSocket 用户集合 全局变量

func init() {
	PersonWSArr = make(map[int64]*ghttp.WebSocket)
}

const RoomDelete = -3 //房间解散
const RoomWait = -2   //房间等待确认中
const GameOver = -1
const DefaultState = 0      //默认不考虑
const StartState = 1        //房间开启
const IdentifySendState = 2 //身份分发完成
const NightState = 3        //入夜（玩家行动） 携带倒计时
const DayState = 4          //白天 携带倒计时
// const RowWheat = 5          //排麦【NowSpeekPlayer 就是现在要开麦说话的人】
// const DivideState = 6       //隔离 有人触发了隔离功能
// const SleepState = 7        //入眠 有人触发了隔离功能
