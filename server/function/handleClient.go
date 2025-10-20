package function

import (
	"bufio"
	"fmt"
	"go_project/src/chat/client/user"
	"net"
	"strings"
)

func HandleClient(conn net.Conn) {

	// 获取用户名
	reader := bufio.NewReader(conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取用户名失败:", err)
		return
	}
	username = strings.TrimSpace(username)
	fmt.Printf("[系统] 用户 %s 连接成功\n", username)
	client := &Client{
		Username: username,
		Conn:     conn,
	}
	lock.Lock()
	users[conn] = client
	lock.Unlock()

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)
		handleClientMessage(client, msg)
	}
	defer func() {
		conn.Close()
		lock.Lock()
		if client, ok := users[conn]; ok {
			fmt.Printf("[系统] 用户 %s 断开连接\n", client.Username)
			delete(users, conn)
			leaveChan <- client.Username
		}
		lock.Unlock()
	}()
}
func handleClientMessage(client *Client, msg string) {
	switch {
	case msg == user.MsgTypeOnline:
		// 用户上线通知
		fmt.Printf("[系统] 用户 %s 进入聊天室\n", client.Username)
		onlineChan <- client.Username
	case msg == user.MsgTypeLeave:
		// 用户离开
		fmt.Printf("[系统] 用户 %s 退出聊天室\n", client.Username)
		lock.Lock()
		delete(users, client.Conn)
		lock.Unlock()
		leaveChan <- client.Username
		client.Conn.Close()
	case msg == user.MsgTypeList:
		// 查看在线用户
		fmt.Printf("[系统] 用户 %s 请求在线用户列表\n", client.Username)
		listChan <- client.Conn
	case strings.HasPrefix(msg, user.MsgTypePrivate+":"):
		// 私聊消息格式: PRIVATE:目标用户:消息内容
		parts := strings.SplitN(strings.TrimPrefix(msg, user.MsgTypePrivate+":"), ":", 2)
		if len(parts) == 2 {
			privateChan <- client.Username + ":" + parts[0] + ":" + parts[1]
		}
	default:
		// 公共群聊消息
		fmt.Printf("[群聊] 用户 %s 发送: %s\n", client.Username, msg)
		msgChan <- fmt.Sprintf("%s: %s", client.Username, msg)
	}
}
