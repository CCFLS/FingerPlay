package httpServer

import (
	"net/http"
	"encoding/json"
	"fmt"
	"strconv"
	"container/list"
	"io"
)

type checkResultParams struct {
	UserToken string `json:"userToken"`
}

type CheckResultData struct {
	Status string `json:"status"`
	Msg string `json:"msg"`
	UserInfos interface{} `json:"userInfos"`
}

type userInfo struct {
	UserToken string `json:"userToken"`
	Winner string `json:"winner"`
}
var status string = "0"

func checkResult(w http.ResponseWriter, req *http.Request)  {
	req.ParseForm()
	parms := req.Form["params"][0]
	checkResultParams := &checkResultParams{}
	json.Unmarshal([]byte(parms),&checkResultParams)
	userToken := checkResultParams.UserToken
	//先判断是不是战局内的人
	isPlayer := false
	for _,value := range Player{
		if value==userToken{
			isPlayer=true
			break
		}
	}
	returnJson := &ReturnJson{}
	if isPlayer{
		fmt.Println("是游戏玩家")
		//判断状态0:皆未出拳1:一方出拳2:胜负已分
		fmt.Println("出拳人数"+strconv.Itoa(len(ActionMap)))
		checkResult := &CheckResultData{}
		if len(ActionMap)==2{//map里有两个人的出拳信息
			checkResult.Status="2"
			checkResult.Msg="胜负已分"
			l := list.New()
			if status!="1"{//不等于1说明需要去计算结果
				status="1"
				l = calResult()
				status="0"
			}
			fmt.Println(l)
			checkResult.UserInfos=l
		}else if len(ActionMap)==1{
			checkResult.Status="1"
			checkResult.Msg="等待对方出拳"
		}else{
			checkResult.Status="0"
			checkResult.Msg="双方均未出拳"
		}
		returnJson.Error = "0"
		returnJson.Msg = "请求成功"
		returnJson.Data = checkResult
	}else {
		returnJson.Error="1"
		returnJson.Msg="该用户未加入战局,查询状态失败"
	}
	jsonStr,_:=json.Marshal(returnJson)
	io.WriteString(w,string(jsonStr))
}


func calResult() *list.List{
	//取当前用户的出拳当基准动作 baseAction 对手是action
	i := 0
	var baseToken string
	var baseAction int
	var token string
	var action int

	for key,value := range ActionMap{
		if i==0 {
			baseToken = key
			baseAction,_ = strconv.Atoi(value)
		}
		if i==1{
			token =key
			action,_ = strconv.Atoi(value)
		}
		i++
	}
	l := list.New()
	if baseAction == 0 {
		if action == 0{
			userInfo1 := &userInfo{}
			userInfo1.UserToken=baseToken
			userInfo1.Winner="T"
			userInfo2 := &userInfo{}
			userInfo2.UserToken=token
			userInfo2.Winner="T"
			l.PushBack(userInfo1)
			l.PushBack(userInfo2)
		}
		if action == 1{
			userInfo1 := &userInfo{}
			userInfo1.UserToken=baseToken
			userInfo1.Winner="F"
			userInfo2 := &userInfo{}
			userInfo2.UserToken=token
			userInfo2.Winner="T"
			l.PushBack(userInfo1)
			l.PushBack(userInfo2)
		}
		if action == 2{
			userInfo1 := &userInfo{}
			userInfo1.UserToken=baseToken
			userInfo1.Winner="T"
			userInfo2 := &userInfo{}
			userInfo2.UserToken=token
			userInfo2.Winner="F"
			l.PushBack(userInfo1)
			l.PushBack(userInfo2)
		}
	}else {
		if baseAction < action {
			userInfo1 := &userInfo{}
			userInfo1.UserToken=baseToken
			userInfo1.Winner="F"
			userInfo2 := &userInfo{}
			userInfo2.UserToken=token
			userInfo2.Winner="T"
			l.PushBack(userInfo1)
			l.PushBack(userInfo2)
		}
		if baseAction > action{
			userInfo1 := &userInfo{}
			userInfo1.UserToken=baseToken
			userInfo1.Winner="T"
			userInfo2 := &userInfo{}
			userInfo2.UserToken=token
			userInfo2.Winner="F"
			l.PushBack(userInfo1)
			l.PushBack(userInfo2)
		}
		if baseAction == action{
			userInfo1 := &userInfo{}
			userInfo1.UserToken=baseToken
			userInfo1.Winner="T"
			userInfo2 := &userInfo{}
			userInfo2.UserToken=token
			userInfo2.Winner="T"
			l.PushBack(userInfo1)
			l.PushBack(userInfo2)
		}
	}
	return l
}
