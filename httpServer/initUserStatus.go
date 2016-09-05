package httpServer

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"
	"io"
	"FingerPlay/redispool"
	"github.com/garyburd/redigo/redis"
	"strings"
)

type initUserReturn struct {
	Role string `json:"role"`
	UserName string `json:"userName"`
	UserHP string `json:"userHP,omitempty"`
}

type initUserStatusParams struct {
	UserName string `json:"userName"`
	Type string `json:"type"`
}

func initUserStatus(w http.ResponseWriter, req *http.Request)  {
	rc := redispool.RedisClient.Get()
	defer rc.Close()
	req.ParseForm()
	parms := req.Form["params"][0]
	initParams := &initUserStatusParams{}
	json.Unmarshal([]byte(parms),&initParams)
	userName := initParams.UserName
	Type := initParams.Type
	if Type =="0"{
		l,_ := redis.Values(rc.Do("KEYS","*"))
		playerNum := 0
		for _,value :=range l{
			if strings.Contains(value,"-player") {
				playerNum++
			}
		}
		fmt.Println("当前战局人数战局人数:"+strconv.Itoa(playerNum))
		if playerNum >=2 {
			value, _ := redis.String(rc.Do("GET", userName+"_player"))
			if value == nil {	//说明不是已经存在的游戏玩家
				Type ="1"	//所以只能成为观战的人
			}
		}
	}
	jsonStr := generatorUser(userName,Type,rc)
	io.WriteString(w,string(jsonStr))
}


func generatorUser(userName string,Type string,rc redis.Conn) string {
	initReturn := &initUserReturn{}
	initReturn.UserName=userName
	if Type == "0"{//参战
		initReturn.Role="0"
		initReturn.UserHP="10"
		userJson,_ := json.Marshal(initReturn)
		rc.Do("SET",userName+"_player",userJson)
	}else if Type == "1"{//观战
		initReturn.Role="1"
		userJson,_ := json.Marshal(initReturn)
		rc.Do("SET",userName+"_watcher",userJson)
	}
	returnJson := &ReturnJson{}
	returnJson.Msg="请求成功"
	returnJson.Error="0"
	returnJson.Data=initReturn
	jsonStr,_:=json.Marshal(returnJson)
	return string(jsonStr)
}
