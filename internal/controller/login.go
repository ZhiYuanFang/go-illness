package controller

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "hello/api/v1"
	"hello/internal/dao"
	"hello/internal/jwt"
	"hello/internal/model/do"
	"hello/internal/model/entity"
	"io/ioutil"
	"net/http"
)

type cLogin struct{}

var Login = cLogin{}

var APPID = "wx070559b32e91d90e"
var APPSECRET = "0f3afea8e65bd34f9db9d6b70652281a"

type WXBody struct {
	openid      string
	session_key string
	unionid     string
	errcode     int
	errmsg      string
}

func (l *cLogin) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	wxCode := req.WXCode
	fmt.Println(wxCode)
	code2sessionHttp := fmt.Sprintf("%s%s%s", "https://api.weixin.qq.com/sns/jscode2session?appid="+APPID+"&secret="+APPSECRET+"&js_code=", wxCode, "&grant_type=authorization_code")
	fmt.Println("请求：", code2sessionHttp)
	wxres, err := http.Get(code2sessionHttp)
	if err != nil {
		return nil, err
	}

	defer wxres.Body.Close()
	wxbodyBytes, _ := ioutil.ReadAll(wxres.Body)

	fmt.Println("返回结果：", string(wxbodyBytes))
	//测试使用
	// wxBody := WXBody{
	// 	openid: wxCode, //todo 测试
	// }
	var wxBody map[string]interface{}
	err = json.Unmarshal(wxbodyBytes, &wxBody)
	if err != nil {
		return nil, err
	}
	openID := wxBody["openid"]
	fmt.Println("当前登录的微信用户openid：", openID)
	//判断并插入新用户
	eUser := entity.User{}
	err = dao.User.Ctx(ctx).Where(do.User{
		OpenId: openID,
	}).Scan(&eUser)

	var id int64
	if err != nil {
		//用户不存在
		id, err = dao.User.Ctx(ctx).InsertAndGetId(do.User{
			OpenId: openID,
		})
		if err != nil {
			return nil, err
		}
	} else {
		//用户存在
		id = eUser.Id
	}

	//执行token获取操作 实现登陆
	accessToken, refreshToken, err := jwt.NewToken(id)
	if err != nil {
		return nil, err
	}
	//返回前端数据
	return &v1.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, err

}
