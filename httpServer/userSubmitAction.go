package httpServer

import (
	"io"
	"encoding/json"
	"net/http"
	"FingerPlay/redispool"
	"strings"
	"github.com/garyburd/redigo/redis"
)

type submitActionParams struct {
	UserName string
	Action string
}

func submitAction(w http.ResponseWriter, req *http.Request){
	rc := redispool.RedisClient.Get()
	defer rc.Close()
	req.ParseForm()
	parms := req.Form["params"][0]
	submitActionParams := &submitActionParams{}
	json.Unmarshal([]byte(parms),&submitActionParams)
	userName := submitActionParams.UserName
	action := submitActionParams.Action
	returnJson := &ReturnJson{}
	//先判断是不是战局内的人
	isPlayer := false
	l,_ := redis.Values(rc.Do("KEYS","*"))
	for _,value :=range l{
		if userName+"player" == value{
			isPlayer = false
			break
		}
	}
	if isPlayer{
		//将出拳放入map,如果已经出拳且未出结果时,将丢弃重复无用的出拳
		ActionMap[userName]=action
		returnJson.Msg="请求成功,已出拳,请等待结果"
		returnJson.Error="0"
	}else {
		returnJson.Msg="该用户未加入战局,出拳失败"
		returnJson.Error="1"
	}
	jsonStr,_:=json.Marshal(returnJson)
	io.WriteString(w,string(jsonStr))
}
