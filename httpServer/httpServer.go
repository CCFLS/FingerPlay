package httpServer

import (
	"net/http"
	"encoding/json"
	"github.com/golibs/uuid"
	"io"
	"sync"
	"fmt"
	"strconv"
	"container/list"
)

var player map[string]string = make(map[string]string)
var watcher map[string]string = make(map[string]string)

type ReturnJson struct {
	Error string `json:"error"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

type initUserReturn struct {
	Role string `json:"role"`
	UserName string `json:"userName"`
	UserToken string `json:"userToken"`
	UserHP string `json:"userHP,omitempty"`
}

type initUserStatusParams struct {
	UserName string `json:"userName"`
	Type string `json:"type"`
}

var lock sync.RWMutex

func initUserStatus(w http.ResponseWriter, req *http.Request)  {
	req.ParseForm()
	parms := req.Form["params"][0]
	initParams := &initUserStatusParams{}
	json.Unmarshal([]byte(parms),&initParams)
	userName := initParams.UserName
	Type := initParams.Type
	uuid := uuid.Rand().Hex()
	if Type =="0"{
		size := len(player) //查看战局人数
		fmt.Println("当前战局人数战局人数:"+strconv.Itoa(len(player)))
		if size >=2 {
			if _, ok := player[userName]; !ok {
				//当前战局人数大于2且该用户不在战局内 则只能加入观战模式
				Type ="1"
			}
		}
	}
	lock.Lock()
	jsonStr := generatorUser(userName,uuid,Type)
	lock.Unlock()
	io.WriteString(w,string(jsonStr))
}

func generatorUser(userName string,uuid string,Type string) string {
	initReturn := &initUserReturn{}
	initReturn.UserName=userName
	if Type == "0"{//参战
		if _, ok := player[userName]; !ok {
			//该用户不存在
			player[userName]=uuid
		}
		initReturn.UserToken=player[userName]
		initReturn.Role="0"
		initReturn.UserHP="10"
	}else if Type == "1"{//观战
		if _, ok := watcher[userName]; !ok {
			//该用户不存在
			watcher[userName]=uuid
		}
		initReturn.UserToken=watcher[userName]
		initReturn.Role="1"
	}
	returnJson := &ReturnJson{}
	returnJson.Msg="请求成功"
	returnJson.Error="0"
	returnJson.Data=initReturn
	jsonStr,_:=json.Marshal(returnJson)
	return string(jsonStr)
}

type submitActionParams struct {
	UserToken string
	Action string
}

var actionMap map[string]string = make(map[string]string)

func submitAction(w http.ResponseWriter, req *http.Request){
	req.ParseForm()
	parms := req.Form["params"][0]
	submitActionParams := &submitActionParams{}
	json.Unmarshal([]byte(parms),&submitActionParams)
	userToken := submitActionParams.UserToken
	action := submitActionParams.Action
	returnJson := &ReturnJson{}
	//先判断是不是战局内的人
	isPlayer := false
	for _,value := range player{
		if value==userToken{
			isPlayer=true
			break
		}
	}
	if isPlayer{
		//将出拳放入map,如果已经出拳且未出结果时,将丢弃重复无用的出拳
		actionMap[userToken]=action
		returnJson.Msg="请求成功,已出拳,请等待结果"
		returnJson.Error="0"
	}else {
		returnJson.Msg="该用户未加入战局,出拳失败"
		returnJson.Error="1"
	}
	jsonStr,_:=json.Marshal(returnJson)
	io.WriteString(w,string(jsonStr))
}

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
	for _,value := range player{
		if value==userToken{
			isPlayer=true
			break
		}
	}
	returnJson := &ReturnJson{}
	if isPlayer{
		fmt.Println("是游戏玩家")
		//判断状态0:皆未出拳1:一方出拳2:胜负已分
		fmt.Println("出拳人数"+strconv.Itoa(len(actionMap)))
		checkResult := &CheckResultData{}
		if len(actionMap)==2{//map里有两个人的出拳信息
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
		}else if len(actionMap)==1{
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

	for key,value := range actionMap{
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

func HttpServerStart()  {
	http.HandleFunc("/initUserStatus",initUserStatus)
	http.HandleFunc("/submitAction",submitAction)
	http.HandleFunc("/checkResult",checkResult)
	http.ListenAndServe("localhost:8080",nil)
}
