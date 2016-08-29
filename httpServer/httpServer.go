package httpServer

import (
	"net/http"
	"encoding/json"
	"github.com/golibs/uuid"
	"io"
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

func initUserStatus(w http.ResponseWriter, req *http.Request)  {
	req.ParseForm()
	parms := req.Form["params"][0]
	initParams := &initUserStatusParams{}
	json.Unmarshal([]byte(parms),&initParams)
	userName := initParams.UserName
	Type := initParams.Type
	uuid := uuid.Rand().Hex()
	if Type == "1"{
		watcher[userName]=uuid
		initReturn := &initUserReturn{}
		initReturn.Role="1"
		initReturn.UserName=userName
		initReturn.UserToken=uuid
		returnJson := &ReturnJson{}
		returnJson.Msg="请求成功"
		returnJson.Error="0"
		returnJson.Data=initReturn
		jsonStr,_:=json.Marshal(returnJson)
		io.WriteString(w,string(jsonStr))
	}else if Type == "0"{
		if len(player)<2{
			player[userName]=uuid
			initReturn := &initUserReturn{}
			initReturn.Role="0"
			initReturn.UserHP="10"
			initReturn.UserName=userName
			initReturn.UserToken=uuid
			returnJson := &ReturnJson{}
			returnJson.Msg="请求成功"
			returnJson.Error="0"
			returnJson.Data=initReturn
			jsonStr,_:=json.Marshal(returnJson)
			io.WriteString(w,string(jsonStr))
		}else {
			watcher[userName]=uuid
			initReturn := &initUserReturn{}
			initReturn.Role="1"
			initReturn.UserName=userName
			initReturn.UserToken=uuid
			returnJson := &ReturnJson{}
			returnJson.Msg="请求成功"
			returnJson.Error="0"
			returnJson.Data=initReturn
			jsonStr,_:=json.Marshal(returnJson)
			io.WriteString(w,string(jsonStr))
		}
	}
}

func HttpServerStart()  {
	http.HandleFunc("/initUserStatus",initUserStatus)
	http.ListenAndServe("localhost:8080",nil)
}
