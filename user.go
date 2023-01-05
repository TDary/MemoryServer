package main

import (
	"MemoryServer/Logs"
	"bytes"
	"encoding/binary"
	"net"
)

type User struct {
	Name    string      //昵称，默认与Add相同
	Addr    string      //地址
	Channel chan string //消息管道
	conn    net.Conn    //连接
	server  *Server     //缓存Server的引用
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
		server:  server,
	}

	//启动协程，监听Channel管道消息
	go user.ListenMessage()

	return user
}

func (us *User) Online() {
	// 用户上线，将用户加入到OnlineMap中，注意加锁操作
	us.server.mapLock.Lock()
	us.server.OnlineMap[us.Name] = us
	us.server.mapLock.Unlock()

	// 广播当前用户上线消息
	us.server.BroadCast(us, "上线啦O(∩_∩)O")
	Logs.Loggers().Print("user Online")
}

func (us *User) Offline() {
	// 用户下线，将用户从OnlineMap中删除，注意加锁
	us.server.mapLock.Lock()
	delete(us.server.OnlineMap, us.Name)
	us.server.mapLock.Unlock()

	// 广播当前用户下线消息
	us.server.BroadCast(us, "下线了o(╥﹏╥)o")
	Logs.Loggers().Print("user Offline")
}

func (us *User) DoMessage(buf []byte, len int) {
	//提取用户的消息(去除'\n')
	msg := string(buf[:len-1])
	Logs.Loggers().Println("DoMessage: ", msg)
	// 调用Server的BroadCast方法
	us.server.BroadCast(us, msg)
}

func (us *User) ListenMessage() {
	for {
		msg := <-us.Channel
		Logs.Loggers().Println("Send msg to client: ", msg, ", len: ", int16(len(msg)))
		bytebuf := bytes.NewBuffer([]byte{})
		// 前两个字节写入消息长度
		binary.Write(bytebuf, binary.BigEndian, int16(len(msg)))
		// 写入消息数据
		binary.Write(bytebuf, binary.BigEndian, []byte(msg))
		// 发送消息给客户端
		us.conn.Write(bytebuf.Bytes())
	}
}
