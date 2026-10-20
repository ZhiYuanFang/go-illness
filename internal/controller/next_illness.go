package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	v1 "hello/api/v1"
	"hello/internal/consts"
	"hello/internal/model/illness"
)

type cNextIllness struct{}

var NextIllness = cNextIllness{}

// 房间集合
var roomArr map[int64]*illness.RoomInfo

// 锁控制
var rxMutex sync.RWMutex

var curRoomId = 0

func init() {
	//初始化房间集合
	roomArr = make(map[int64]*illness.RoomInfo)
}

// 第一步 玩家开启天梯 / 教学局（IsTutorial=true 时固定 4 人且非 lying）
func (c *cNextIllness) CreateRoom(ctx context.Context, req *v1.IllnessCreateRoomReq) (res *v1.IllnessCreateRoomRes, err error) {
	//判断当前玩家是否已经在游戏房间内
	playerId := req.Id
	roomId := GetPlayerRoomID(playerId)
	if roomId > 0 {
		//已经在房间
		go SendWS2PlayerWithRoomInfo(playerId, roomId)
	} else {
		isLying := req.IsLyingType
		roomNumber := req.RoomNumber
		isTutorial := req.IsTutorial
		if isTutorial {
			// 教学局：服务端权威人数与模式
			isLying = false
			roomNumber = 4
		}
		_, err = CreateNewRoom(ctx, playerId, isLying, roomNumber, isTutorial)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// SkipTutorial 主动跳过教学并持久化完成标记（无需已在教学房内）。
func (c *cNextIllness) SkipTutorial(ctx context.Context, req *v1.IllnessSkipTutorialReq) (res *v1.IllnessSkipTutorialRes, err error) {
	if req.Id <= 0 {
		return nil, errors.New("无效用户")
	}
	if err = MarkTutorialCompleted(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.IllnessSkipTutorialRes{}, nil
}

// 第1.1步 玩家继续本房天梯
func (c *cNextIllness) ContinueGame(ctx context.Context, req *v1.IllnessContinueRoomReq) (res *v1.IllnessContinueRoomRes, err error) {

	roomId := req.RoomId
	room, ok := GetRoom(roomId)
	if ok {
		//获取玩家的RANK值
		player, err := GetPlayerInfo(ctx, req.Id)
		player.IsRoomManager = true
		if err == nil {
			room.PlayerArr[req.Id] = player
			rank := player.WantLimitRoomNumber
			room.LimitRoomNumber = rank
			room.NumNight = 0
			//获取本房间的活人人数

			var persons map[int64]*illness.PlayerInfo
			persons = make(map[int64]*illness.PlayerInfo)
			for k, v := range room.PlayerArr {
				if k > 0 {
					if v.IsRoomManager {
						v.HasPrepared = false
					}
					v.IsQuarantine = false
					v.InfectedSeri = 0
					v.IsInfected = false
					v.Identity = 0
					v.CheckContent = nil
					v.TagIdentify = 0
					v.TagInfo = ""

					persons[k] = v
				}
			}

			if len(persons) == 0 && room.IsLyingType {
				return nil, errors.New("当前房间人数不匹配，请重新创建房间")
			}

			if len(persons) > rank {
				err = errors.New(fmt.Sprint("当前天梯", rank, " 当前人数超额，请减员"))
			} else {
				//重置人数
				room.PlayerArr = persons
				fmt.Println("rank", rank)
				fmt.Println("personLen:", len(persons))
				var dis = rank - len(persons)
				for i := 0; i < dis; i++ {
					machinID := int64(i+1) * -1
					machinPlayer := CreateMachinInfo(machinID)
					room.PlayerArr[machinID] = machinPlayer
					fmt.Println("add:", machinID)
				}
				fmt.Println("size:", len(room.PlayerArr))
				//角色集合
				room.AllIdentifyArr = GetIdentifyArrForCurNumber(rank)
				room.GameStat = consts.RoomWait

				//通知玩家房间重置完毕
				go SendWS2AllPlayerOnlyRoom(room.RoomId)
			}
		}
	} else {
		errors.New("本房间已经解散")
	}
	return nil, err
}

// 第二步 玩家移除机器人
func (c *cNextIllness) RemoveMachin(ctx context.Context, req *v1.IllnessRemoveMachinReq) (res *v1.IllnessRemoveMachinRes, err error) {
	room, ok := GetRoom(req.RoomId)
	if ok {
		if !room.PlayerArr[req.Id].IsRoomManager {
			return nil, errors.New("只有房主可以踢人")
		}
		var tempPlayer = room.PlayerArr[req.PlayerId]
		if tempPlayer != nil {
			delete(room.PlayerArr, tempPlayer.Id)
			//告知前端刷新
			go SendWS2AllPlayerOnlyRoom(req.RoomId)
		} else {
			err = errors.New("当前人员不在此房间内")
		}
	} else {
		err = errors.New("当前房间不存在")
	}
	return nil, err
}

// 第三步 加入房间
func (c *cNextIllness) JoinRoom(ctx context.Context, req *v1.IllnessJoinRoomReq) (res *v1.IllnessJoinRoomRes, err error) {
	room, ok := GetRoom(req.RoomId)
	if ok {
		if len(room.PlayerArr) < room.LimitRoomNumber {
			//房间未满员
			player, err := GetPlayerInfo(ctx, req.PlayerId)
			if err == nil {
				room.PlayerArr[req.PlayerId] = player
				go SendWS2AllPlayerOnlyRoom(req.RoomId)
			}
		} else {
			err = errors.New("房间已满员")
		}
	} else {
		err = errors.New("当前房间不存在")
	}
	return nil, err
}

// 第四步 准备 判断所有人准备 则开始
func (c *cNextIllness) PreparePlay(ctx context.Context, req *v1.IllnessPlayerPrepareReq) (res *v1.IllnessPlayerPrepareRes, err error) {
	room, ok := GetRoom(req.RoomId)
	if ok {
		if err = TutorialGate(room, TutorialActionPrepare); err != nil {
			return nil, err
		}

		var tempPlayer = room.PlayerArr[req.Id]
		if tempPlayer != nil {
			tempPlayer.HasPrepared = true

			//判断是否全员准备
			isAllReady := true
			if len(room.PlayerArr) == room.LimitRoomNumber {
				personNum := 0 //活人数量
				for _, pi := range room.PlayerArr {
					if !pi.HasPrepared {
						isAllReady = false
					}
					if pi.Id > 0 {
						personNum++
					}
				}

				if personNum == 1 && room.IsLyingType {
					return nil, errors.New("当前房间类型不匹配，请重新创建房间")
				}
			} else {
				isAllReady = false
			}
			if isAllReady {
				AdvanceTutorial(room)
				go StartGame(room)
			} else {
				//当前玩家已准备，告知前端刷新
				go SendWS2AllPlayerOnlyRoom(req.RoomId)
			}
		} else {
			err = errors.New("当前人员不在此房间内")
		}
	} else {
		err = errors.New("当前房间不存在")
	}
	return nil, err
}

// 第五步 进入黑夜
// 契约：感染仅在夜间流程由服务端施加（无开局预感染）。
func ContinueNight(room *illness.RoomInfo) {

	fmt.Println("第五步 进入黑夜")
	//记录当前为第几夜
	room.NumNight = room.NumNight + 1

	//通知进入天黑
	room.CurrentNightTimeStamp = time.Now().Unix() //设置天黑的初始时间
	go SendWS2AllPlayerOnlyStatus(room.RoomId, consts.NightState)

	noInfectedIdArr := []int64{}
	//已感染的数量
	hasInfectedNumber := 0
	//筛选未被感染的数组
	for k, pi := range room.PlayerArr {
		if !pi.IsInfected {
			noInfectedIdArr = append(noInfectedIdArr, k)
		} else {
			hasInfectedNumber++
		}

		//重置当夜查询的对象ID
		pi.TonightVerPlayerId = 0
	}

	// 夜间随机感染一名未感染者（规格：无开局预感染）
	if len(noInfectedIdArr) > 0 {
		id := noInfectedIdArr[rand.Intn(len(noInfectedIdArr))]
		player := room.PlayerArr[id]
		player.IsInfected = true
		fmt.Println("感染新玩家=>", id)
		hasInfectedNumber++
		player.InfectedSeri = hasInfectedNumber //从1开始

		//今晚感染之后 所有人都被感染了
		if hasInfectedNumber >= room.LimitRoomNumber {
			time.Sleep(3 * time.Second)
			//休眠3秒后通知 游戏结束
			go GameOver(room, false)
			return
		}

		//获取未隔离的数组
		tempInQuarantineArr := []*illness.PlayerInfo{}
		for _, pi := range room.PlayerArr {
			if !pi.IsQuarantine {
				tempInQuarantineArr = append(tempInQuarantineArr, pi)
			}
		}

		//找到机器人，如果是医生 让他从未隔离的玩家中 随机验证一名玩家
		for k, pi := range room.PlayerArr {
			if k < 0 {
				if pi.Identity < 3 && len(room.PlayerJSONArr) > 0 {
					tempPlayer := room.PlayerArr[tempInQuarantineArr[int64(rand.Intn(len(tempInQuarantineArr)))].Id]
					verPlayerInfected(room, pi, tempPlayer)
				}
			}
		}

		//休眠一定时间后 进入白天
		time.Sleep(time.Duration(room.NightTime) * time.Second)
		go SendWS2AllPlayerOnlyStatus(room.RoomId, consts.DayState)
		fmt.Println("进入白天")
	} else {
		go GameOver(room, false)
	}
}

// 第六步 白天的时候送入隔离
func (c *cNextIllness) DividePerson(ctx context.Context, req *v1.IllnessDividePersonReq) (res *v1.IllnessDividePersonRes, err error) {
	room, ok := roomArr[req.RoomId]
	if ok {
		if err = TutorialGate(room, TutorialActionQuarantine); err != nil {
			return nil, err
		}

		for k, pi := range room.PlayerArr {
			fmt.Println(k, "是否感染：", pi.IsInfected)
		}

		var tempPlayer = room.PlayerArr[req.PlayerId]
		if tempPlayer != nil {
			//判断当前人员是否被感染
			tempPlayer.IsQuarantine = true
			AdvanceTutorial(room)
			if tempPlayer.IsInfected {
				//隔离正确
				//通知房间内所有玩家，有人送入隔离了
				go SendWS2AllPlayerOnlyRoom(req.RoomId)

				//判断是否游戏胜利
				go func(room *illness.RoomInfo) {
					end := true
					for _, pi := range room.PlayerArr {
						if pi.IsInfected && !pi.IsQuarantine {

							fmt.Println("未找到", pi.Id)
							end = false
						}
					}
					if end {

						fmt.Println("3")
						GameOver(room, true)
					}
				}(room)
			} else {
				//隔离失败，游戏结束（误隔离健康玩家）
				go GameOver(room, false)
			}
		} else {
			err = errors.New("当前人员不在此房间内")
		}
	} else {
		err = errors.New("当前房间不存在")
	}
	return nil, err
}

// 第七步 白天的时候 进入黑夜
func (c *cNextIllness) GoNight(ctx context.Context, req *v1.IllnessPersonSleepReq) (res *v1.IllnessPersonSleepRes, err error) {
	room, ok := roomArr[req.RoomId]
	if ok {
		if err = TutorialGate(room, TutorialActionSleep); err != nil {
			return nil, err
		}
		fmt.Println("手动进入黑夜")
		AdvanceTutorial(room)
		go ContinueNight(room)
	} else {
		err = errors.New("当前房间不存在")
	}
	return nil, err
}

// 添加机器
func (c *cNextIllness) AddMachin(ctx context.Context, req *v1.IllnessAddMachinReq) (res *v1.IllnessAddMachinRes, err error) {

	room, ok := GetRoom(req.RoomId)

	if ok {
		oldMachinIds := []int64{}
		for k, _ := range room.PlayerArr {
			if k < 0 {
				oldMachinIds = append(oldMachinIds, k)
			}
		}

		newID := int64(0)
		for {
			newID--
			isContain := false
			for _, v := range oldMachinIds {
				if v == newID {
					isContain = true
				}
			}
			if !isContain {
				break
			}
		}

		room.PlayerArr[newID] = CreateMachinInfo(newID)

		//判断是否满足开始游戏的条件
		//判断是否全员准备
		isAllReady := true
		if len(room.PlayerArr) == room.LimitRoomNumber {
			for _, pi := range room.PlayerArr {
				if !pi.HasPrepared {
					isAllReady = false
				}
			}
		} else {
			isAllReady = false
		}
		if isAllReady {
			go StartGame(room)
		} else {
			//当前机器人已添加，告知前端刷新
			go SendWS2AllPlayerOnlyRoom(req.RoomId)
		}
	} else {
		err = errors.New("当前房间不存在")
	}
	return nil, nil
}

// 强制退出【扣除经验，并结束本局】
func (c *cNextIllness) ExitRoom(ctx context.Context, req *v1.IllnessExitRoomReq) (res *v1.IllnessExitRoomRes, err error) {
	room, ok := GetRoom(req.RoomId)
	if ok {
		var tempPlayer = room.PlayerArr[req.Id]
		if tempPlayer != nil {
			if tempPlayer.IsRoomManager {
				go RoomEnd(room.RoomId)
			} else {
				go func(room *illness.RoomInfo) {
					GameOver(room, false)
					delete(room.PlayerArr, tempPlayer.Id)
				}(room)

			}
		} else {
			err = errors.New("当前人员不在此房间内")
		}
	} else {
		err = errors.New("当前房间不存在")
	}
	return nil, err
}

// 标记身份
func (c *cNextIllness) TagIdentify(ctx context.Context, req *v1.IllnessTagIdentifyReq) (res *v1.IllnessTagIdentifyRes, err error) {
	room, ok := roomArr[req.RoomId]
	if ok {
		if err = TutorialGate(room, TutorialActionTag); err != nil {
			return nil, err
		}
		var tempPlayer = room.PlayerArr[req.PlayerId]
		if tempPlayer != nil {
			tempPlayer.TagIdentify = req.TagIndentify
			AdvanceTutorial(room)
			go SendWS2AllPlayerOnlyRoom(req.RoomId)
		} else {
			err = errors.New("当前人员不在此房间内")
		}
	} else {
		err = errors.New("当前房间不存在")
	}
	return nil, err
}

// 标记自定义
func (c *cNextIllness) TagInfo(ctx context.Context, req *v1.IllnessTagInfoReq) (res *v1.IllnessTagInfoRes, err error) {
	room, ok := roomArr[req.RoomId]
	if ok {
		if err = TutorialGate(room, TutorialActionTag); err != nil {
			return nil, err
		}
		var tempPlayer = room.PlayerArr[req.PlayerId]
		if tempPlayer != nil {
			tempPlayer.TagInfo = req.TagInfo
			AdvanceTutorial(room)
			go SendWS2AllPlayerOnlyRoom(req.RoomId)
		} else {
			err = errors.New("当前人员不在此房间内")
		}
	} else {
		err = errors.New("当前房间不存在")
	}
	return nil, err
}

// -----------------------------------------------------------
// 获取机器人身份信息
func CreateMachinInfo(machinID int64) (playerInfo *illness.PlayerInfo) {
	return &illness.PlayerInfo{
		Id:            machinID,
		NickName:      fmt.Sprint(machinID*-1, "号"),
		Avatar:        "https://120.55.50.105:9686/images/robot_avatar.webp",
		Lv:            0,
		StartWaitTime: time.Now().Unix(),
		IsSleep:       true, //默认进入睡眠状态
		HasPrepared:   true,
	}
}

// 获取玩家所在房间ID
func GetPlayerRoomID(playerId int64) (roomId int64) {
	rxMutex.Lock()
	defer rxMutex.Unlock()
	for k, ri := range roomArr {
		for k2, _ := range ri.PlayerArr {
			if k2 == playerId {
				//当前玩家已经在房间内
				return k
			}
		}
	}
	return 0
}

// 为玩家开辟房间；isTutorial 时强制 4 人非 lying，并初始化教学 tip。
func CreateNewRoom(ctx context.Context, playerId int64, isLyingType bool, roomNumber int, isTutorial bool) (room *illness.RoomInfo, err error) {
	playArr := make(map[int64]*illness.PlayerInfo)
	//获取玩家信息
	comeInPlayerInfo, err := GetPlayerInfo(ctx, playerId)
	if err != nil {
		return nil, err
	}
	comeInPlayerInfo.IsRoomManager = true
	playArr[comeInPlayerInfo.Id] = comeInPlayerInfo
	if isTutorial {
		isLyingType = false
		roomNumber = 4
	} else if !isLyingType {
		roomNumber = comeInPlayerInfo.WantLimitRoomNumber
	}
	if roomNumber < 3 {
		roomNumber = 3
	}
	if !isLyingType {
		//生成机器人信息
		fmt.Println("获取机器人信息")
		for i := 0; i < roomNumber-1; i++ {
			machinID := int64(i+1) * -1
			machinPlayer := CreateMachinInfo(machinID)
			playArr[machinID] = machinPlayer
		}
	}

	//角色集合
	allIdentifyArr := GetIdentifyArrForCurNumber(roomNumber)

	//申请房间ID
	curRoomId++
	room = &illness.RoomInfo{
		IsLyingType:     isLyingType,      //是否说谎模式
		RoomId:          int64(curRoomId), //房间ID
		PlayerArr:       playArr,          //房间内成员列表
		AllIdentifyArr:  allIdentifyArr,   //当前房间的所有身份集合
		NightTime:       8,                //晚上默认时长10秒  考虑limitRoomNumber 定黑夜时间比例
		GameStat:        consts.RoomWait,  //默认状态
		LimitRoomNumber: roomNumber,       //房间人数
		IsTutorial:      isTutorial,
		TutorialStep:    0,
	}
	if isTutorial {
		ApplyTutorialTip(room)
	}

	//将房间加入房间集合
	PutRoomIntoArr(room)

	//通知玩家房间创建完毕
	go SendWS2AllPlayerOnlyStatus(room.RoomId, consts.RoomWait)

	return room, nil
}

// 获取玩家信息
func GetPlayerInfo(ctx context.Context, playerId int64) (playerInfo *illness.PlayerInfo, err error) {
	userInfoRes, err := User.UserInfo(ctx, &v1.UserInfoReq{
		Id: playerId,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("获取最新RANK", userInfoRes.Rank)
	WantLimitRoomNumber := userInfoRes.Rank
	if WantLimitRoomNumber < 3 {
		WantLimitRoomNumber = 3
	}
	return &illness.PlayerInfo{
		Id:                  userInfoRes.Id,
		NickName:            userInfoRes.NickName,
		Avatar:              userInfoRes.Avatar,
		Lv:                  userInfoRes.Lv,
		StartWaitTime:       time.Now().Unix(),
		HasPrepared:         false,
		WantLimitRoomNumber: WantLimitRoomNumber,
	}, nil
}

// 获取角色集合
func GetIdentifyArrForCurNumber(limitRoomNumber int) (allIdentifyArr []int) {
	allIdentifyArr = []int{1, 2, 3}
	badNumber := 1
	for i := 0; i < limitRoomNumber-3; i++ {
		if i%2 > 0 {
			//6人 2庸医 12人 3庸医 18人 4庸医
			bsix := int(limitRoomNumber / 6)
			if badNumber < bsix+1 {
				allIdentifyArr = append(allIdentifyArr, 1)
				badNumber = badNumber + 1
			} else {
				allIdentifyArr = append(allIdentifyArr, 2)
			}
		} else {
			allIdentifyArr = append(allIdentifyArr, 3)
		}
	}
	return allIdentifyArr
}

// 将房间加入集合
func PutRoomIntoArr(room *illness.RoomInfo) {
	rxMutex.Lock()
	defer rxMutex.Unlock()
	roomArr[room.RoomId] = room

}

// 获取房间
func GetRoom(roomID int64) (room *illness.RoomInfo, ok bool) {
	rxMutex.Lock()
	defer rxMutex.Unlock()
	room, ok = roomArr[roomID]
	return room, ok
}

// 游戏开始
func StartGame(room *illness.RoomInfo) {
	//全员已准备 开始进入身份分发
	_, err := IdentityDistribution(room.RoomId, room.AllIdentifyArr)
	if err != nil {
		//身份分发失败， 房间解散
		go RoomEnd(room.RoomId)
	} else {
		time.Sleep(5 * time.Second)
		//5秒后 进入黑夜（在黑夜里 执行操作 并一定时间后进入白天）
		go ContinueNight(room)
	}
}

// 房间解散
func RoomEnd(roomId int64) {
	SendWS2AllPlayerOnlyStatus(roomId, consts.RoomDelete)
	rxMutex.Lock()
	defer rxMutex.Unlock()
	delete(roomArr, roomId)
}

// 身份分发
func IdentityDistribution(roomID int64, candidateArr []int) (playerMap map[int64]*illness.PlayerInfo, err error) {
	//身份分配
	//方案一：1个庸医 5个医生 6个平民， 一半医生 一半平民， 医生里随机一个是庸医 1，2，3身份牌
	//候选身份初始集合的复制，因为GO是没有弱引用的，全部都是强引用，这里一定要复制
	candidateArrCopy := make([]int, len(candidateArr))
	copy(candidateArrCopy, candidateArr)
	room, ok := GetRoom(roomID)
	if !ok {
		return nil, errors.New("当前房间不存在")
	}
	playerArr := room.PlayerArr
	for playerId, pi := range playerArr {
		endPlayer, identifyIndex, err := IdentityDistribution_SingleIdentity(pi, candidateArrCopy)
		if err != nil {
			//分配身份出现问题，抛出异常
			break
		}

		//当前玩家身份分配完毕
		playerArr[playerId] = endPlayer
		//从身份初始复制集合中删除当前身份，以供下一次循环不要有冲突身份...go语言中没有直接从数组删除的功能，用拼接数组的方式实现
		if identifyIndex >= len(candidateArrCopy) {
			candidateArrCopy = candidateArrCopy[:identifyIndex]
		} else {
			candidateArrCopy = append(candidateArrCopy[:identifyIndex], candidateArrCopy[(identifyIndex+1):]...)
		}
	}

	//分发身份通知到对应玩家 游戏进程为2
	go SendWS2AllPlayerOnlyStatus(roomID, consts.IdentifySendState)
	return playerArr, err
}

// 身份分发-单人身份随机 waitIdentity 尚未分配完毕的身份
func IdentityDistribution_SingleIdentity(inPlayer *illness.PlayerInfo, waitIdentity []int) (outPlayer *illness.PlayerInfo, identifyIndex int, err error) {
	if len(waitIdentity) > 0 {
		identifyIndex = rand.Intn(len(waitIdentity))
		inPlayer.Identity = waitIdentity[identifyIndex]
		inPlayer.IsSleep = true //默认进入睡眠状态
		// fmt.Println(fmt.Sprint("单人分发随机身份：", inPlayer.Identity))
		return inPlayer, identifyIndex, nil
	} else {
		return nil, 0, nil
	}
}

// 游戏结束
// 天梯契约：非 lying 且非教学局时，房主 rank 胜 +1 / 负 −1 / 下限 3（失败不会重置为 3）。
func GameOver(room *illness.RoomInfo, IsSuccess bool) {
	room.IsSuccess = IsSuccess
	// 非说谎、非教学：修改房主天梯值
	if !room.IsLyingType && !room.IsTutorial {
		for k, pi := range room.PlayerArr {
			if pi.IsRoomManager {

				rank := pi.WantLimitRoomNumber
				fmt.Println("pi.WantLimitRoomNumber", rank)
				if IsSuccess {
					rank++
				} else {
					rank--
					if rank < 3 {
						rank = 3
					}
				}
				fmt.Println("修改RANK值", rank)
				User.UpDateUserRank(context.Background(), &v1.UserRankUpdateReq{
					Id:   k,
					Rank: rank,
				})

				break
			}
		}
	}
	// 教学局胜利：持久化教学完成标记
	if room.IsTutorial && IsSuccess {
		for k, pi := range room.PlayerArr {
			if pi.IsRoomManager && k > 0 {
				_ = MarkTutorialCompleted(context.Background(), k)
				break
			}
		}
	}

	SendWS2AllPlayerOnlyStatus(room.RoomId, consts.GameOver)
}

func verPlayerInfected(room *illness.RoomInfo, curPlayer *illness.PlayerInfo, tempPlayer *illness.PlayerInfo) (isInfected bool, err error) {
	if (curPlayer.TonightVerPlayerId != 0 && curPlayer.TonightVerPlayerId == tempPlayer.Id) || curPlayer.TonightVerPlayerId == 0 {
		isRepeatCheck := curPlayer.TonightVerPlayerId != 0
		if tempPlayer != nil {
			curPlayer.TonightVerPlayerId = tempPlayer.Id

			var isInfected = tempPlayer.IsInfected
			if !room.IsLyingType || curPlayer.Id < 0 { //不是说谎模式下\或者是机器人 庸医获得相反结果
				if curPlayer.Identity != 2 { //如果不是医生，得到相反的验证结果
					isInfected = !isInfected
				}
			}

			//记录查验结果
			var str = "未感染"
			if isInfected {
				str = "已感染"
			}
			if !isRepeatCheck {
				endStr := fmt.Sprint(tempPlayer.NickName, " ", str)
				// fmt.Println(endStr)
				curPlayer.CheckContent = append(curPlayer.CheckContent, endStr)
			}

			return isInfected, nil
		} else {
			err = errors.New("当前人员不在此房间内")
		}
	} else {
		err = errors.New("今夜已经验证过其它人")
	}
	return false, err
}

// 医生验证感染信息
func (c *cNextIllness) VerPersonInfected(ctx context.Context, req *v1.IllnessVerInfectedPersonReq) (res *v1.IllnessVerInfectedPersonRes, err error) {
	room, ok := roomArr[req.RoomId]
	if ok {
		if err = TutorialGate(room, TutorialActionVerify); err != nil {
			return nil, err
		}
		var tempPlayer = room.PlayerArr[req.PlayerId]
		var curPlayer = room.PlayerArr[req.Id]
		isInfected, err := verPlayerInfected(room, curPlayer, tempPlayer)
		if err == nil {
			AdvanceTutorial(room)
			go SendWS2AllPlayerOnlyRoom(req.RoomId)
		}
		return &v1.IllnessVerInfectedPersonRes{
			IsInfected: isInfected,
		}, err
	} else {
		err = errors.New("当前房间不存在")
	}
	return nil, err
}

//------------------------------------------------------------------------

// 发送通知到房间内的所有玩家
func SendWS2AllPlayerOnlyStatus(roomId int64, gameStat int) (err error) {
	room, ok := GetRoom(roomId)
	if ok {
		fmt.Println(fmt.Sprint("发送通知:", gameStat, room.UdpPort))
		if room != nil {
			room.GameStat = gameStat
			fmt.Println("修改GameStat", room.GameStat)
			SendWS2AllPlayerOnlyRoom(roomId)
		} else {
			fmt.Println("当前房间不存在")
		}
	}
	return nil
}

// 发送通知到房间内的所有玩家
func SendWS2AllPlayerOnlyRoom(roomId int64) (err error) {
	if room, _ := GetRoom(roomId); room != nil {
		playerArr := room.PlayerArr
		if playerArr != nil {
			for k, _ := range playerArr {
				_, ok := consts.PersonWSArr[k]
				if ok {
					SendWS2PlayerWithRoomInfo(k, roomId)
				} else {
					// fmt.Println("找不到用户", k, "的长连接")
				}
			}
		} else {
			fmt.Println("当前房间没有玩家")
		}
	} else {
		fmt.Println("当前房间不存在")
	}
	return err
}

func SendWS2PlayerWithRoomInfoJustOnlyPlayer(playerId int64) (err error) {
	if len(roomArr) > 0 {
		for k, ri := range roomArr {
			for k2, _ := range ri.PlayerArr {
				if k2 == playerId {
					go SendWS2PlayerWithRoomInfo(playerId, k)
					return nil
				}
			}
		}
	}
	return errors.New("当前玩家尚未匹配成功")
}

// 发送给固定玩家房间信息(用于玩家重连)
func SendWS2PlayerWithRoomInfo(playerId int64, roomId int64) (err error) {
	fmt.Println(fmt.Sprint("发送给固定玩家房间信息:", roomId))
	if room, _ := GetRoom(roomId); room != nil {
		fmt.Println("`````````>", room)
		ws, ok := consts.PersonWSArr[playerId]
		if ok {
			playerArr := room.PlayerArr
			//规范返回玩家列表信息
			curPlayerArr := make([]*illness.PlayerInfo, len(playerArr))
			tempKeys := []int{}
			for ki, _ := range playerArr { //玩家列表里 有玩家的身份和ID 便于中间数据丢失的存在
				tempKeys = append(tempKeys, int(ki))
			}
			// sort.Ints(tempKeys)
			sort.Slice(tempKeys, func(i, j int) bool {
				return tempKeys[i] > tempKeys[j]
			})
			for index := 0; index < len(tempKeys); index++ {
				pi := playerArr[int64(tempKeys[index])]
				curPlayerArr[index] = &illness.PlayerInfo{
					Id:                pi.Id,
					NickName:          pi.NickName,
					Avatar:            pi.Avatar,
					IsQuarantine:      pi.IsQuarantine,
					InfectedSeri:      pi.InfectedSeri,
					IsInfected:        pi.IsInfected,
					IsSleep:           pi.IsSleep,
					BeReportedNum:     pi.BeReportedNum,
					ReportPlayerIdArr: pi.ReportPlayerIdArr,
					Identity:          pi.Identity,
					HasPrepared:       pi.HasPrepared,
					IsRoomManager:     pi.IsRoomManager,
					TagIdentify:       pi.TagIdentify,
					TagInfo:           pi.TagInfo,
				}

				if room.IsLyingType { //说谎游戏的时候， 庸医需要自行判断是否说谎
					if room.GameStat == consts.GameOver {
						curPlayerArr[index].CheckContent = pi.CheckContent
						if pi.Identity == 1 { //庸医的获胜条件与正常相反
							room.IsSuccess = !room.IsSuccess
						}
					} else {
						//游戏未结束，验证信息 如果是机器人 则告知玩家，如果是玩家且是本人 则告知
						if pi.Id > 0 && playerId == pi.Id {
							curPlayerArr[index].CheckContent = pi.CheckContent
						} else if pi.Id < 0 {
							curPlayerArr[index].CheckContent = pi.CheckContent
						}
					}
				} else {
					//不说谎模式下，庸医身份隐藏
					if pi.Identity == 1 && room.GameStat != consts.GameOver {
						curPlayerArr[index].Identity = 2
					}
					//验证信息显示
					curPlayerArr[index].CheckContent = pi.CheckContent
				}
			}

			room.PlayerJSONArr = curPlayerArr //房间玩家列表
			if ws != nil {
				if err = SendWS2PlayerWithWs(ws, room); err != nil {
					fmt.Println(fmt.Sprintf("发送给用户", playerId, " 失败"))
				}
			} else {
				//某用户的长连接断开了，环境模拟：匹配成功的瞬间，用户的设备出现问题，长连接断开，需要用户自己去判断 并重新链接。 后期可以加上消息推送，告知用户
			}
		} else {
			fmt.Println("找不到用户", playerId, "的长连接")
		}
	} else {
		fmt.Println("当前房间不存在")
	}

	return err
}

// 发送通知到固定玩家
func SendWS2PlayerWithWs(ws *ghttp.WebSocket, msg any) (err error) {
	if ws == nil {
		return errors.New("长连接为空")
	}
	msgBytes, err := json.Marshal(msg) //使用已经被格式化的房间信息作为输出
	if err == nil {
		err = ws.WriteMessage(1, msgBytes)
	} else {
		//序列化房间信息错误
		fmt.Println("序列化数据信息错误")
	}
	return err
}
