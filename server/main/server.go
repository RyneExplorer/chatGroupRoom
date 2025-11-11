package main

import (
	"fmt"
	"go_project/src/chat/server/service"
	"net"
)

func main() {
	fmt.Println("服务器启动，监听端口 15000...")
	listener, err := net.Listen("tcp", "127.0.0.1:15000")
	if err != nil {
		fmt.Println("监听失败:", err)
		return
	}
	defer listener.Close()

	go service.HeartBeatCheck()
	go service.BroadcastLoop()
	go service.BroadcastFromStreamLoop()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受连接失败:", err)
			continue
		}
		go service.ChatRoomEntry(conn)
	}
}
