package httpServer

import (
	"net/http"
	"encoding/json"
	"io"
	"strings"
	"FingerPlay/redispool"
	"github.com/garyburd/redigo/redis"
)

type userQuitParams struct {
	UserName string
}

func userQuit(w http.ResponseWriter, req *http.Request)  {
	rc := redispool.RedisClient.Get()
	defer rc.Close()
	req.ParseForm()
	parms := req.Form["params"][0]
	userQuitParams := &userQuitParams{}
	json.Unmarshal([]byte(parms),&userQuitParams)
	userName := userQuitParams.UserName
	rc.Do("DELETE",userName+"_player")
	rc.Do("DELETE",userName+"_watcher")
	returnJson := &ReturnJson{}
	returnJson.Error="OK"
	returnJson.Msg="退出成功"
	jsonStr,_:=json.Marshal(returnJson)
	io.WriteString(w,string(jsonStr))
}

