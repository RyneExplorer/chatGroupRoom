package service

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func ChatRoomEntry(conn net.Conn) {
	defer conn.Close()
	err := ResetAllUserStatus(LEAVE_STATUS)
	if err != nil {
		log.Fatalf("重置用户状态失败: %v", err)
	}

	var username string
	reader := bufio.NewReader(conn)
	// 处理注册/登录请求
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("客户端 %s 在认证阶段断开: %v", conn.RemoteAddr().String(), err)
			return
		}
		msg = strings.TrimSpace(msg)

		// 获取认证用户名
		username, err = handleAuthMessage(conn, msg)
		if err == nil {
			break
		} else {
			log.Printf("用户 %s 认证时出现错误: %v", username, err)
			continue
		}
	}
	// 认证成功后处理聊天消息
	err = handleChatMessages(conn, username)
	if err != nil {
		log.Fatalf("服务端出现错误: %v", err)
	}

}
