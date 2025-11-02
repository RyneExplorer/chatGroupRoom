package service

import (
	"bufio"
	"errors"
	"fmt"
	"go_project/src/chat/client/user"
	"log"
	"math"
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
			sendResponse(conn, "注册失败! 服务器异常...")
			return username, err
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
			sendResponse(conn, "登录失败! 服务器异常...")
			return username, err
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

func handleChatInfo(conn net.Conn, username string) {
	client := &Client{
		Username:       username,
		Conn:           conn,
		LastActiveTime: time.Now(),
	}

	Lock.Lock()
	Users[conn] = client
	Lock.Unlock()

	fmt.Printf("[系统] 用户 %s 连接成功\n", username)
	level, cur, next := GetUserLevelAndProgress(getUserLevel(client.Username), int(GetScore(client.Username)))
	info := "\n-----------------------\n"
	info += fmt.Sprintf("[系统] 欢迎回来! %s\n你当前活跃度等级为Lv.%d (%d/%d)\n", client.Username, level, cur, next)
	conn.Write([]byte(info))
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		Lock.Lock()
		client.LastActiveTime = time.Now()
		Lock.Unlock()
		msg = strings.TrimSpace(msg)
		handleClientMessage(client, msg)
	}

	conn.Close()
	Lock.Lock()
	// 处理客户端异常断开的问题 1.释放user  2.更新用户状态
	if client, ok := Users[conn]; ok {
		log.Printf("[系统] 用户 %s 异常断开连接\n", client.Username)
		delete(Users, conn)
		LeaveChan <- client.Username
		err := UpdateStatus(client.Username, LEAVE_STATUS)
		if err != nil {
			log.Printf("警告: 更新异常下线用户状态失败: %v", err)
		}
	}
	Lock.Unlock()
}

func handleClientMessage(client *Client, msg string) {
	switch {
	case msg == user.MsgTypeOnline:
		// 用户上线通知
		fmt.Printf("[系统] 用户 %s 进入聊天室\n", client.Username)
		OnlineChan <- client.Username
	case msg == user.MsgTypeLeave:
		// 用户离开
		fmt.Printf("[系统] 用户 %s 退出\n", client.Username)
		Lock.Lock()
		delete(Users, client.Conn)
		LeaveChan <- client.Username
		err := UpdateStatus(client.Username, LEAVE_STATUS)
		if err != nil {
			log.Printf("更新下线用户 %s 状态失败: %v", client.Username, err)
		}
		Lock.Unlock()
		client.Conn.Close()
	case msg == user.MsgTypeList:
		// 查看在线用户
		fmt.Printf("[系统] 用户 %s 请求在线用户列表\n", client.Username)
		ListChan <- client.Conn
	case strings.HasPrefix(msg, user.MsgTypePrivate+":"):
		// 私聊消息格式: PRIVATE:目标用户:消息内容
		// strings.SplitN返回一个切割n-1次且含有n个元素的字符串切片
		parts := strings.SplitN(strings.TrimPrefix(msg, user.MsgTypePrivate+":"), ":", 2)
		if len(parts) == 2 {
			PrivateChan <- client.Username + ":" + parts[0] + ":" + parts[1]
		}
	case msg == user.MsgTypeRank:
		// 获取活跃度排名
		fmt.Printf("[系统] 用户 %s 请求活跃度排行\n", client.Username)
		err := showRank(client)
		if err != nil {
			log.Printf("%v", err)
			client.Conn.Write([]byte("获取活跃度排名失败! 请稍后重试...\n"))
		}
	default:
		// 公共群聊消息
		fmt.Printf("[群聊] 用户 %s 发送: %s\n", client.Username, msg)
		MsgChan <- fmt.Sprintf("%s: %s", client.Username, msg)
		CheckUserLevel(client)
	}
}
func showRank(client *Client) error {
	data, err := GetRankData(client)
	if err != nil {
		return err
	}
	// 设定用户名最大显示宽度，不足补空格
	const NameWidth = 12
	rank := "\n------活跃度排行榜------\n"
	for i, val := range data {
		// %-*s: 左对齐，* 表示宽度由参数传入
		level := getUserLevel(val)
		rank += fmt.Sprintf("%d. %-*s Lv%d\n", i+1, NameWidth, val, level)
	}
	client.Conn.Write([]byte(rank + "\n"))
	return nil
}

// ScoreExpToLevel 等级增长机制
func ScoreExpToLevel(level int) (score int) {
	// base乘以level的exponent的幂次方 = redis中的score并返回
	base := 5.0     // 基础值
	exponent := 1.5 // 增长指数
	score = int(math.Pow(float64(level), exponent) * base)
	return score
}

// GetNewLevel 获取用户活跃度新等级
func GetNewLevel(level, score int) int {
	for ScoreExpToLevel(level) <= score {
		level++
	}
	return level
}

// GetUserLevelAndProgress 检查当前用户等级和经验
func GetUserLevelAndProgress(level, score int) (int, int, int) {
	if score < 5 {
		return 1, score, 5
	}
	for ScoreExpToLevel(level) <= score {
		level++
	}
	currentExp := score
	nextExp := ScoreExpToLevel(level)
	return level, currentExp, nextExp
}

func CheckUserLevel(client *Client) {
	// 增加用户活跃度
	err := AddUserActivity(client)
	if err != nil {
		log.Printf("%v", err)
	}
	// 判断用户是否升级
	oldLevel := getUserLevel(client.Username)
	score := GetScore(client.Username)
	newLevel := GetNewLevel(oldLevel, int(score))
	if newLevel > oldLevel {
		msg := fmt.Sprintf("[系统] 恭喜用户 %s 等级提升! Lv.%d --> Lv.%d", client.Username, oldLevel, newLevel)
		LevelChan <- msg
		err = UpdateUserLevel(client.Username, newLevel)
		if err != nil {
			log.Printf("%v", err)
		} else {
			client.Level = newLevel
		}
	}
}
