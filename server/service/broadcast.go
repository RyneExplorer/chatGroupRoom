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
		case msg := <-MsgChan:
			broadcastToAll(fmt.Sprintf("[群聊] %s", msg))
		case online := <-OnlineChan:
			broadcastToAll(fmt.Sprintf("[系统] 用户 %s 加入了聊天室", online))
		case leave := <-LeaveChan:
			broadcastToAll(fmt.Sprintf("[系统] 用户 %s 离开了聊天室", leave))
		case level := <-LevelChan:
			broadcastToAll(fmt.Sprintf(level))
		case private := <-PrivateChan:
			// 私聊消息直接发送给目标用户
			parts := strings.SplitN(private, ":", 3)
			if len(parts) == 3 {
				fmt.Printf("[私聊] %s -> %s: %s\n", parts[0], parts[1], parts[2])
				sendPrivateMessage(parts[0], parts[1], parts[2])
			}
		case conn := <-ListChan:
			if _, ok := Users[conn]; ok {
				sendUserList(conn)
			}
		}
	}
}

func broadcastToAll(msg string) {
	Lock.Lock()
	defer Lock.Unlock()

	for conn := range Users {
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			log.Printf("广播发送失败: %v", err)
			return
		}
	}
}

// 发送私聊消息
func sendPrivateMessage(from, to, content string) {
	Lock.Lock()
	defer Lock.Unlock()

	found := false
	for conn, client := range Users {
		if client.Username == to {
			found = true
			_, err := conn.Write([]byte(fmt.Sprintf("[私聊] %s对你私聊: %s\n", from, content)))
			if err != nil {
				log.Printf("用户 %s 发送私聊失败: %v", from, err)
				return
			}
		}
	}

	// 如果目标用户不在线，通知发送者
	if !found {
		for conn, client := range Users {
			if client.Username == from {
				conn.Write([]byte(fmt.Sprintf("[系统] 用户 %s 离线, 发送失败!\n"+"\n", to)))
				break
			}
		}
	}
}

// 发送在线用户列表
func sendUserList(conn net.Conn) {
	Lock.Lock()
	defer Lock.Unlock()

	list := "\n----- 在线用户 -----\n"
	for _, client := range Users {
		list += fmt.Sprintf("- %s\n", client.Username)
	}
	conn.Write([]byte(list + "\n"))
}
