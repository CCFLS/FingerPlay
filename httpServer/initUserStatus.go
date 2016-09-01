package httpServer

import (
	"encoding/json"
	"sync"
	"net/http"
	"github.com/golibs/uuid"
	"fmt"
	"strconv"
	"io"
)

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
	lock.Lock()
	if Type =="0"{
		size := len(Player) //查看战局人数
		fmt.Println("当前战局人数战局人数:"+strconv.Itoa(len(Player)))
		if size >=2 {
			if _, ok := Player[userName]; !ok {
				//当前战局人数大于2且该用户不在战局内 则只能加入观战模式
				Type ="1"
			}
		}
	}
	jsonStr := generatorUser(userName,uuid,Type)
	lock.Unlock()
	io.WriteString(w,string(jsonStr))
}

func generatorUser(userName string,uuid string,Type string) string {
	initReturn := &initUserReturn{}
	initReturn.UserName=userName
	if Type == "0"{//参战
		if _, ok := Player[userName]; !ok {
			//该用户不存在
			Player[userName]=uuid
		}
		initReturn.UserToken=Player[userName]
		initReturn.Role="0"
		initReturn.UserHP="10"
	}else if Type == "1"{//观战
		if _, ok := Watcher[userName]; !ok {
			//该用户不存在
			Watcher[userName]=uuid
		}
		initReturn.UserToken=Watcher[userName]
		initReturn.Role="1"
	}
	returnJson := &ReturnJson{}
	returnJson.Msg="请求成功"
	returnJson.Error="0"
	returnJson.Data=initReturn
	jsonStr,_:=json.Marshal(returnJson)
	return string(jsonStr)
}
