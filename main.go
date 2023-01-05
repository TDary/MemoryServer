//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
package main

import (
	"MemoryServer/Logs"
	"time"
)

func main() {
	//启动服务器
	go StartServer()

	//写入其他业务逻辑
	Logs.Loggers().Print("这是一个Go服务端，socket消息广播功能")

	//防止主线程退出
	for {
		time.Sleep(1 * time.Second)
	}
}

func StartServer() {
	server := NewServer("127.0.0.1", 8231)
	server.Start()
}
