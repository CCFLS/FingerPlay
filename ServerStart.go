package main

import (
	"fmt"
	"FingerPlay/httpServer"
	"FingerPlay/redispool"
	"github.com/garyburd/redigo/redis"
)

func main() {
	fmt.Println("服务已启动")
	httpServer.HttpServerStart()
}

func testredis(){
	redispool.Init()
	rc := redispool.RedisClient.Get()
	defer rc.Close()
	_, err :=  rc.Do("SET","test","hhhh")
	if err != nil {
		fmt.Println(err)
		return
	}
	value, err := redis.String(rc.Do("GET", "test"))
	if err != nil {
		fmt.Println("fail")
	}
	fmt.Println(value)
}
func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}