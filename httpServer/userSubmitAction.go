package httpServer

import (
	"io"
	"encoding/json"
	"net/http"
)

type submitActionParams struct {
	UserToken string
	Action string
}

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
	for _,value := range Player{
		if value==userToken{
			isPlayer=true
			break
		}
	}
	if isPlayer{
		//将出拳放入map,如果已经出拳且未出结果时,将丢弃重复无用的出拳
		ActionMap[userToken]=action
		returnJson.Msg="请求成功,已出拳,请等待结果"
		returnJson.Error="0"
	}else {
		returnJson.Msg="该用户未加入战局,出拳失败"
		returnJson.Error="1"
	}
	jsonStr,_:=json.Marshal(returnJson)
	io.WriteString(w,string(jsonStr))
}
