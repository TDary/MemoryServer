//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
package main

import (
	"MemoryServer/HttpServer"
	"MemoryServer/Logs"
	"MemoryServer/SocketServer"
	"time"
)

func main() {
	//启动与Unity通信服务器
	go SocketServer.StartServer()

	//启动与用户通信服务器
	go HttpServer.ListenAndServer("10.11.144.31:9070")

	//写入其他业务逻辑
	Logs.Loggers().Print("服务器启动成功")

	//防止主线程退出
	for {
		time.Sleep(1 * time.Second)
	}
}
