package httpServer

import (
	"net/http"
	"encoding/json"
	"FingerPlay/redispool"
	"github.com/garyburd/redigo/redis"
	"strings"
	"io"
)

type checkTheOtherSideParams struct {
	UserName string `json:"userName"`
}

type hasOtherReturn struct {
	Role string `json:"role"`
	UserName string `json:"userName"`
	UserHP string `json:"userHP,omitempty"`
}

func checkTheOtherSide(w http.ResponseWriter, req *http.Request)  {
	rc := redispool.RedisClient.Get()
	defer rc.Close()
	req.ParseForm()
	parms := req.Form["params"][0]
	otherSideParams := &checkTheOtherSideParams{}
	json.Unmarshal([]byte(parms),&otherSideParams)
	userName := otherSideParams.UserName
	l,_ := redis.Values(rc.Do("KEYS","*"))
	initUserJson := ""
	for _,value :=range l{
		if strings.Contains(value,"-player") {
			if userName+"_player" != value{
				initUserJson,_ = rc.Do("GET",value)
			}
		}
	}
	io.WriteString(w,"{\"error\":\"OK\",\"msg\":\"0\",\"data\":"+initUserJson+"}")
}