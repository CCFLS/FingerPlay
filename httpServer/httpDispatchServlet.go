package httpServer

import "net/http"

//当前加入战局的人数。猜拳游戏,同一局最多两人
var Player map[string]string = make(map[string]string)
//观战的人。人数无限制
var Watcher map[string]string = make(map[string]string)
//保存出拳信息。双方出拳后,可以计算结果,计算完成后清空数据
var ActionMap map[string]string = make(map[string]string)

type ReturnJson struct {
	Error string `json:"error"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func HttpServerStart()  {
	http.HandleFunc("/initUserStatus",initUserStatus)
	http.HandleFunc("/checkTheOtherSide",checkTheOtherSide)
	http.HandleFunc("/submitAction",submitAction)
	http.HandleFunc("/checkResult",checkResult)
	http.HandleFunc("/UserQuit",userQuit)
	http.ListenAndServe("0.0.0.0:8080",nil)
}
