package service

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func BroadcastLoop() {
	for {
		select {
		case msg := <-msgChan:
			broadcastToAll(fmt.Sprintf("[群聊] %s", msg))
		case online := <-onlineChan:
			broadcastToAll(fmt.Sprintf("[系统] 用户 %s 加入了聊天室", online))
		case leave := <-leaveChan:
			broadcastToAll(fmt.Sprintf("[系统] 用户 %s 离开了聊天室", leave))
		case private := <-privateChan:
			// 私聊消息直接发送给目标用户
			parts := strings.SplitN(private, ":", 3)
			if len(parts) == 3 {
				fmt.Printf("[私聊] %s -> %s: %s\n", parts[0], parts[1], parts[2])
				sendPrivateMessage(parts[0], parts[1], parts[2])
			}
		case conn := <-listChan:
			if _, ok := users[conn]; ok {
				sendUserList(conn)
			}
		}
	}
}

func broadcastToAll(msg string) {
	lock.Lock()
	defer lock.Unlock()

	for conn := range users {
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			log.Fatal("发送失败:", err)
		}
	}
}

// 发送私聊消息
func sendPrivateMessage(from, to, content string) {
	lock.Lock()
	defer lock.Unlock()

	found := false
	for conn, client := range users {
		if client.Username == to {
			found = true
			_, err := conn.Write([]byte(fmt.Sprintf("[私聊] %s对你私聊: %s\n", from, content)))
			if err != nil {
				log.Fatal("发送客户端私聊失败: ", err)
			}
		}
	}

	// 如果目标用户不在线，通知发送者
	if !found {
		for conn, client := range users {
			if client.Username == from {
				conn.Write([]byte(fmt.Sprintf("[系统] 用户 %s 离线, 发送失败!\n", to)))
				break
			}
		}
	}
}

// 发送在线用户列表
func sendUserList(conn net.Conn) {
	lock.Lock()
	defer lock.Unlock()

	list := "\n----- 在线用户 -----\n"
	for _, client := range users {
		list += fmt.Sprintf("- %s\n", client.Username)
	}
	list += "\n"
	conn.Write([]byte(list))
}
