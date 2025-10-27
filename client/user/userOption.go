package user

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func Verify(conn net.Conn) {
	for {
		err := RegisterLogin(conn)
		if err == nil {
			break
		}
		fmt.Println("认证失败:", err)
		fmt.Println("请重新尝试...")
	}

}
func ChatMenuLoop(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("----- 聊天室菜单 -----")
		fmt.Println("1. 进入聊天室")
		fmt.Println("2. 退出")
		fmt.Println("3. 私聊用户")
		fmt.Println("4. 查看在线用户")
		fmt.Println("请选择操作: ")

		if scanner.Scan() {
			op := strings.TrimSpace(scanner.Text())
			switch op {
			case "1":
				sendMessage(conn, MsgTypeOnline)
				chatLoop(conn, scanner)
			case "2":
				sendMessage(conn, MsgTypeLeave)
				return
			case "3":
				fmt.Print("请输入对方用户名: ")
				if !scanner.Scan() {
					fmt.Println("获取用户名失败!")
					return
				}
				target := strings.TrimSpace(scanner.Text())
				fmt.Print("请输入私聊内容: ")
				if !scanner.Scan() {
					fmt.Println("读取消息失败!")
					return
				}
				content := strings.TrimSpace(scanner.Text())
				sendPrivateMessage(conn, target, content)
				time.Sleep(time.Millisecond * 10)
			case "4":
				sendMessage(conn, MsgTypeList)
				time.Sleep(time.Millisecond * 10)
			default:
				fmt.Println("无效操作，请重新选择!")
			}
		}
	}
}

// 进入聊天室聊天
func chatLoop(conn net.Conn, scanner *bufio.Scanner) {
	fmt.Println("[系统] 你已进入聊天室，输入消息发送群聊，输入'q'返回菜单")
	for {
		if scanner.Scan() {
			msg := strings.TrimSpace(scanner.Text())
			if msg == "q" {
				fmt.Println("[系统] 你已退出聊天室")
				return
			}
			if msg != "" {
				_, err := conn.Write([]byte(msg + "\n"))
				if err != nil {
					fmt.Println("发送失败:", err)
					return
				}
			}
		}
	}
}

func sendMessage(conn net.Conn, msgType string) {
	_, err := conn.Write([]byte(msgType + "\n"))
	if err != nil {
		fmt.Println("发送失败:", err)
	}
}

func sendPrivateMessage(conn net.Conn, target, content string) {
	msg := MsgTypePrivate + ":" + target + ":" + content
	_, err := conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("私聊发送失败:", err)
	}
}
