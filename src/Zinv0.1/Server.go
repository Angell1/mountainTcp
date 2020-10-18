package main

import (
	"../znet"
)

func main() {
	//创建Server句柄
	Server := znet.NerServer()
	Server.Run()
}
