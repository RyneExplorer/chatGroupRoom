package service

import (
	"bufio"
	"errors"
	"fmt"
	"go_project/src/chat/client/user"
	"log"
	"net"
	"strings"
	"time"
)

func handleAuthMessage(conn net.Conn, msg string) (string, error) {
	parts := strings.SplitN(msg, ":", 3)
	if len(parts) != 3 {
		sendResponse(conn, "数据格式错误!")
		return "", ErrDataFormat
	}
	msgType := parts[0]
	username := parts[1]
	password := parts[2]
	switch {
	case msgType == user.MsgTypeRegister:
		err := Register(username, password)
		if err == nil {
			sendResponse(conn, "注册成功!")
			return username, nil
		} else if errors.Is(err, ErrUsernameExists) {
			sendResponse(conn, "用户名已存在, 请重新创建!")
			return username, ErrUsernameExists
		} else {
			sendResponse(conn, "注册失败!")
			return username, fmt.Errorf("注册失败")
		}
	case msgType == user.MsgTypeLogin:
		err := Login(username, password)
		if err == nil {
			sendResponse(conn, "登录成功!")
			return username, nil
		} else if errors.Is(err, ErrUsernameNotExists) {
			sendResponse(conn, "该用户名不存在!")
			return username, ErrUsernameNotExists
		} else if errors.Is(err, ErrInvalidInput) {
			sendResponse(conn, "用户名或密码不能为空!")
			return username, ErrInvalidInput
		} else if errors.Is(err, ErrUsernamePassData) {
			sendResponse(conn, "密码错误!")
			return username, ErrUsernamePassData
		} else if errors.Is(err, ErrUserIsOnline) {
			sendResponse(conn, "该用户已经上线, 无法重复登录!")
			return username, ErrUserIsOnline
		} else {
			sendResponse(conn, "登录失败!")
			return username, fmt.Errorf("登录失败")
		}
	default:
		sendResponse(conn, "该操作不支持!")
		return username, ErrNotSupported
	}
}

// 响应客户端登录注册功能
func sendResponse(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg + "\n"))
	if err != nil {
		log.Printf("服务端响应出错: %v", err)
	}
}

func handleChatMessages(conn net.Conn, username string) error {
	client := &Client{
		Username:       username,
		Conn:           conn,
		LastActiveTime: time.Now(),
	}

	lock.Lock()
	users[conn] = client
	lock.Unlock()

	fmt.Printf("[系统] 用户 %s 连接成功\n", username)
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("读取用户 %s 消息失败: %v", username, err)
			break
		}
		lock.Lock()
		client.LastActiveTime = time.Now()
		lock.Unlock()
		msg = strings.TrimSpace(msg)
		err = handleClientMessage(client, msg)
		if err != nil {
			return err
		}
	}

	conn.Close()
	lock.Lock()
	// 处理客户端异常断开的问题 1.释放user  2.更新用户状态
	if client, ok := users[conn]; ok {
		log.Printf("[系统] 用户 %s 异常断开连接\n", client.Username)
		delete(users, conn)
		leaveChan <- client.Username
		err := UpdateStatus(client.Username, LEAVE_STATUS)
		if err != nil {
			log.Printf("警告: 更新异常下线用户状态失败!")
			return err
		}
	}
	lock.Unlock()
	return nil

}

func handleClientMessage(client *Client, msg string) error {
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
		leaveChan <- client.Username
		err := UpdateStatus(client.Username, LEAVE_STATUS)
		if err != nil {
			return ErrUpdateStatusFailed
		}
		lock.Unlock()
		client.Conn.Close()
	case msg == user.MsgTypeList:
		// 查看在线用户
		fmt.Printf("[系统] 用户 %s 请求在线用户列表\n", client.Username)
		listChan <- client.Conn
	case strings.HasPrefix(msg, user.MsgTypePrivate+":"):
		// 私聊消息格式: PRIVATE:目标用户:消息内容
		// strings.SplitN返回一个切割n-1次且含有n个元素的字符串切片
		parts := strings.SplitN(strings.TrimPrefix(msg, user.MsgTypePrivate+":"), ":", 2)
		if len(parts) == 2 {
			privateChan <- client.Username + ":" + parts[0] + ":" + parts[1]
		}
	default:
		// 公共群聊消息
		fmt.Printf("[群聊] 用户 %s 发送: %s\n", client.Username, msg)
		msgChan <- fmt.Sprintf("%s: %s", client.Username, msg)
	}
	return nil
}
