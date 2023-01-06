package SocketServer

import (
	"MemoryServer/Logs"
	"fmt"
	"io"
	"net"
	"sync"
)

var AllInfo []InfoData

type InfoData struct {
	Ip           string
	UnityVersion string
	IsTakeSimple bool
	FileName     string
}

type Server struct {
	Ip   string
	Port int

	//在线用户容器
	OnlineMap map[string]*User
	//用户列容器锁，对容器进行操作时会加锁
	mapLock sync.RWMutex

	//消息广播的管道
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (serv *Server) Start() {
	//socket监听
	Listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serv.Ip, serv.Port))
	if err != nil {
		Logs.Loggers().Println("net.listen err:", err)
		return
	}

	//程序退出时关闭监听，注意defer关键字用途
	defer Listen.Close()

	//启动一个协程来执行ListenMessager
	go serv.ListenMessager()

	for {
		//Accept，此处会堵塞，当有客户端连接时才会继续往后面执行
		conn, err := Listen.Accept()
		if err != nil {
			Logs.Loggers().Println("listener accept err:", err)
			continue
		}

		//启动一个协程去处理
		go serv.Handler(conn)
	}
}

func (serv *Server) Handler(conn net.Conn) {

	// 构造User对象，NewUser全局方法在user.go脚本中
	user := NewUser(conn, serv)

	// 用户上线
	user.Online()

	// 启动一个协程
	go func() {
		buf := make([]byte, 4096)
		for {
			// 从Conn中读取消息
			length, err := conn.Read(buf)
			if length == 0 {
				// 用户下线
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				Logs.Loggers().Println("Conn Read err:", err)
				return
			}

			// 用户针对msg进行消息处理
			user.DoMessage(buf, length)
		}
	}()
}

func (serv *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]: " + msg

	serv.Message <- sendMsg
}

func (serv *Server) ReceiveFromApi(username string, msg string) {
	sendMsg := "[" + username + "]: " + msg
	serv.Message <- sendMsg
}

func (serv *Server) ListenMessager() {
	for {
		// 从Message管道中读取消息
		msg := <-serv.Message

		// 加锁
		serv.mapLock.Lock()
		// 遍历在线用户，把广播消息同步给在线用户
		for _, user := range serv.OnlineMap {
			// 把要广播的消息写到用户管道中
			user.Channel <- msg
		}
		// 解锁
		serv.mapLock.Unlock()
	}
}
