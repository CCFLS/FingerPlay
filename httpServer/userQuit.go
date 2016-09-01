package httpServer

import (
	"net/http"
	"encoding/json"
	"io"
)

type userQuitParams struct {
	UserToken string
}

func userQuit(w http.ResponseWriter, req *http.Request)  {
	req.ParseForm()
	parms := req.Form["params"][0]
	userQuitParams := &userQuitParams{}
	json.Unmarshal([]byte(parms),&userQuitParams)
	userToken := userQuitParams.UserToken
	for key,value:=range Player{
		if value == userToken{
			delete(Player,key)
		}
	}
	for key,value:=range Watcher{
		if value == userToken{
			delete(Watcher,key)
		}
	}
	returnJson := &ReturnJson{}
	returnJson.Error="OK"
	returnJson.Msg="退出成功"
	jsonStr,_:=json.Marshal(returnJson)
	io.WriteString(w,string(jsonStr))
}

