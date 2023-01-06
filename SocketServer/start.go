package SocketServer

import (
	"strings"
)

func StartServer() {
	server := NewServer("10.11.144.31", 9080)
	go CheckMessage(server)
	server.Start()
}

func CheckMessage(serv *Server) {
	for {
		if len(AllInfo) != 0 {
			for i := 0; i < len(AllInfo); i++ {
				msg := AllInfo[i].UnityVersion + "," + AllInfo[i].FileName
				serv.ReceiveFromApi(AllInfo[i].Ip, msg)
			}
		}
	}
}

func GetData(data string) {
	var current_ip string
	var unityVersion string
	var info InfoData
	current := strings.Split(data, "&")
	for i := 0; i < len(current); i++ {
		if strings.Contains(current[i], "ip") {
			cdata := strings.Split(current[i], "=")
			current_ip = cdata[1]
			info.Ip = current_ip
		} else if strings.Contains(current[i], "unityVersion") {
			cdata := strings.Split(current[i], "=")
			unityVersion = cdata[1]
			info.UnityVersion = unityVersion
		} else if strings.Contains(current[i], "fileName") {
			cdata := strings.Split(current[i], "=")
			filename := cdata[1]
			info.FileName = filename
		}
	}
	if info.Ip != "" && info.UnityVersion != "" {
		info.IsTakeSimple = false
		AllInfo = append(AllInfo, info)
	}
}
