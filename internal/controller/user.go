// 用户信息
package controller

import (
	"context"
	"errors"
	"fmt"
	v1 "hello/api/v1"
	"hello/internal/dao"
	"hello/internal/model/do"
	"hello/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type cUser struct{}

var User = cUser{}

func (u *cUser) UserInfo(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error) {
	eUser := entity.User{}
	err = dao.User.Ctx(ctx).Where(do.User{
		Id: req.Id,
	}).Scan(&eUser)
	if err != nil {
		//用户不存在
		return nil, err
	} else {
		//用户存在
		return &v1.UserInfoRes{
			Id:                eUser.Id,
			NickName:          eUser.NickName,
			Avatar:            eUser.Avatar,
			Lv:                eUser.Lv,
			Rank:              eUser.Rank,
			MaxRank:           eUser.MaxRank,
			UseLiveDay:        eUser.UseLiveDay,
			TutorialCompleted: eUser.TutorialCompleted,
		}, err
	}
}

func (u *cUser) UpDateUserInfo(ctx context.Context, req *v1.UserInfoUpdateReq) (res *v1.UserInfoRes, err error) {

	//获取用户原来的头像昵称
	eUser := entity.User{}
	err = dao.User.Ctx(ctx).Where(do.User{
		Id: req.Id,
	}).Scan(&eUser)
	if err == nil {
		res = &v1.UserInfoRes{
			Id:                eUser.Id,
			NickName:          eUser.NickName,
			Avatar:            eUser.Avatar,
			Lv:                eUser.Lv,
			Rank:              eUser.Rank,
			MaxRank:           eUser.MaxRank,
			UseLiveDay:        eUser.UseLiveDay,
			TutorialCompleted: eUser.TutorialCompleted,
		}

		if req.Avatar != "" {
			res.Avatar = req.Avatar
		}

		if req.NickName != "" {
			res.NickName = req.NickName
		}

		//更新头像和昵称
		_, err := dao.User.Ctx(ctx).Update(g.Map{"Avatar": res.Avatar, "NickName": res.NickName}, "id", res.Id)

		return res, err
	} else {
		return nil, err
	}

}

func (u *cUser) UpDateUserRank(ctx context.Context, req *v1.UserRankUpdateReq) (res *v1.UserInfoRes, err error) {
	eUser := entity.User{}
	err = dao.User.Ctx(ctx).Where(do.User{
		Id: req.Id,
	}).Scan(&eUser)
	if err == nil {
		//更新Rank
		if eUser.MaxRank < req.Rank {
			_, err := dao.User.Ctx(ctx).Update(g.Map{"Rank": req.Rank, "MaxRank": req.Rank}, "id", req.Id)
			return res, err
		} else {
			_, err := dao.User.Ctx(ctx).Update(g.Map{"Rank": req.Rank}, "id", req.Id)
			return res, err
		}

	} else {
		return nil, err
	}

}

// 判断是否可以复活
func (u *cUser) UserReliveJudge(ctx context.Context, req *v1.UserReliveJudgeReq) (res *v1.UserInfoRes, err error) {
	eUser := entity.User{}
	err = dao.User.Ctx(ctx).Where(do.User{
		Id: req.Id,
	}).Scan(&eUser)
	if err == nil {
		//todo 记录获取复活时间
		lastTime := eUser.UseLiveDay
		curTime := gtime.Now()

		fmt.Println("上一次使用复活的时间", lastTime)
		fmt.Println("当前时间", curTime)
		if lastTime.Year() > 1 {
			//有上一次复活的时间
			curYear := curTime.Year()
			curMonth := curTime.Month()
			curDay := curTime.Day()

			lastYear := lastTime.Year()
			lastMonth := lastTime.Month()
			lastDay := lastTime.Day()

			fmt.Println("上一次使用复活的时间", lastYear, lastMonth, lastDay)
			fmt.Println("当前时间", curYear, curMonth, curDay)
			if curYear == lastYear && curMonth == lastMonth && curDay == lastDay {
				//当天已经复活过了
				return nil, errors.New("已经没有复活机会了")
			} else {
				return nil, err
			}
		} else {
			fmt.Println("未曾使用过复活")
			return nil, err
		}

	} else {
		return nil, err
	}

}

// MarkTutorialCompleted 将用户教学完成标记置为 true（胜利结算或 skipTutorial）。
func MarkTutorialCompleted(ctx context.Context, userId int64) error {
	_, err := dao.User.Ctx(ctx).Data(g.Map{
		dao.User.Columns().TutorialCompleted: true,
	}).Where(dao.User.Columns().Id, userId).Update()
	return err
}

// 复活时间更新
func (u *cUser) UserRelive(ctx context.Context, req *v1.UserReliveUpdateReq) (res *v1.UserInfoRes, err error) {
	eUser := entity.User{}
	err = dao.User.Ctx(ctx).Where(do.User{
		Id: req.Id,
	}).Scan(&eUser)
	if err == nil {
		//todo 记录获取复活时间
		curTime := gtime.Now()
		_, err := dao.User.Ctx(ctx).Update(g.Map{"UseLiveDay": curTime}, "id", req.Id)
		return nil, err
	} else {
		return nil, err
	}

}
