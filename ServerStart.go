package main

import (
	"fmt"
	"FingerPlay/httpServer"
)

func main() {
	fmt.Println("服务已启动")
	httpServer.HttpServerStart()
}
