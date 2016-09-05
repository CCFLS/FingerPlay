package httpServer

import (
	"net/http"
	"encoding/json"
)

type checkTheOtherSideParams struct {
	UserName string `json:"userName"`
}

type hasOtherReturn struct {
	Role string `json:"role"`
	UserName string `json:"userName"`
	UserToken string `json:"userToken"`
	UserHP string `json:"userHP,omitempty"`
}

func checkTheOtherSide(w http.ResponseWriter, req *http.Request)  {
	req.ParseForm()
	parms := req.Form["params"][0]
	otherSideParams := &checkTheOtherSideParams{}
	json.Unmarshal([]byte(parms),&otherSideParams)
	userName := otherSideParams.UserName
	hasOther := false
	for key,_:=range Player{
		if key != userName{
			hasOther=true
			userName=key
		}
	}
	if hasOther{
		initReturn := &initUserReturn{}
		initReturn.UserName=userName
		initReturn.UserToken=Player[userName]
		initReturn.Role="0"
		initReturn.UserHP="10"
	}
}