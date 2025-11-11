package service

import (
	"log"
	"net"
)

func ChatRoomEntry(conn net.Conn) {
	defer conn.Close()

	// 处理注册/登录请求
	username, err := authenticateUser(conn)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	// 认证成功后调用handleChatInfo处理聊天消息
	handleChatInfo(conn, username)
}
