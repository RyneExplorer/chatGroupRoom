package main

import (
	"fmt"
	"go_project/src/chat/client/user"
	"net"
)

func main() {
	fmt.Println("客户端启动...")
	conn, err := net.Dial("tcp", "127.0.0.1:15000")
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()

	user.Verify(conn)
	go user.ReadMessage(conn)
	user.ChatMenuLoop(conn)
}
