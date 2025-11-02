package user

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// RegisterLogin 客户端注册与登录入口
func RegisterLogin(conn net.Conn) error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("=====  登录/注册  =====")
		fmt.Println("1.登录")
		fmt.Println("2.注册")
		fmt.Println("3.退出")
		fmt.Println("请选择操作: ")
		if scanner.Scan() {
			op := strings.TrimSpace(scanner.Text())
			switch op {
			case "1":
				err := login(conn, MsgTypeLogin)
				if err != nil {
					log.Printf("程序错误! %v", err)
					continue
				}
				return nil
			case "2":
				err := register(conn, MsgTypeRegister)
				if err != nil {
					log.Printf("程序错误! %v", err)
					continue
				}
				return nil
			case "3":
				fmt.Println("系统退出")
				os.Exit(0)
			default:
				fmt.Println("无效的操作, 请重新输入!")
				fmt.Println()
			}
		}
	}
}

func register(conn net.Conn, msgType string) error {
	reader := bufio.NewReader(os.Stdin)
	var username, password string
	for {
		fmt.Println("请输入用户名(3-12位): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("读取用户名失败: %w", err)
		}
		username = strings.TrimSpace(input)

		if len(username) >= 3 && len(username) <= 12 {
			break
		}
		fmt.Println("用户名长度必须在3-12位之间，请重新输入！")
		continue
	}

	for {
		fmt.Println("请输入密码(6-18位): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("读取密码失败: %w", err)
		}
		password = strings.TrimSpace(input)
		if len(password) >= 6 && len(password) <= 18 {
			break
		}
		fmt.Println("密码长度必须在6-18位之间，请重新输入！")
		continue
	}

	msg := fmt.Sprintf("%s:%s:%s", msgType, username, password)
	_, err := conn.Write([]byte(msg + "\n"))
	if err != nil {
		return fmt.Errorf("发送注册信息失败: %w", err)
	}

	// 响应服务器结果
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("服务器响应出错! 请稍后重试... %w", err)
	}
	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "注册成功!") {
		fmt.Println("[系统] 注册成功! 即将进入聊天室...")
		fmt.Println()
		return nil
	} else {
		fmt.Println("[系统] 注册失败: ", response+"\n")
		return fmt.Errorf("%s", response)
	}
}

func login(conn net.Conn, msgType string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("请输入用户名: ")
	usernameInput, _ := reader.ReadString('\n')
	username := strings.TrimSpace(usernameInput)

	fmt.Println("请输入密码: ")
	passwordInput, _ := reader.ReadString('\n')
	password := strings.TrimSpace(passwordInput)

	// 发送登录请求
	msg := fmt.Sprintf("%s:%s:%s", msgType, username, password)
	_, err := conn.Write([]byte(msg + "\n"))
	if err != nil {
		return fmt.Errorf("发送登录信息失败: %w", err)
	}

	// 响应服务端
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("服务器响应出错! 请稍后重试... %w", err)
	}
	response = strings.TrimSpace(response)

	if strings.HasPrefix(response, "登录成功!") {
		fmt.Println("[系统] 登录成功! 即将进入聊天室...")
		fmt.Println()
		return nil
	} else {
		fmt.Println("[系统] 登录失败: ", response+"\n")
		return fmt.Errorf("%s", response)
	}
}
