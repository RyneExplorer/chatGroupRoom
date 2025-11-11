package user

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func ReadMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		// 先读取消息长度
		lengthBytes := make([]byte, 4)
		_, err := io.ReadFull(reader, lengthBytes)
		if err != nil {
			log.Fatalf("程序错误! 服务器连接断开: %v", err)
		}

		// 将字节数组转换为整数，表示消息的长度
		msgLength := int(binary.BigEndian.Uint32(lengthBytes))

		// 根据长度读取消息内容
		msgBytes := make([]byte, msgLength)
		_, err = io.ReadFull(reader, msgBytes)
		if err != nil {
			log.Printf("程序错误! 读取消息失败: %v", err)
		}

		msg := string(msgBytes)
		fmt.Println(strings.TrimSpace(msg))
	}
}

//func ReadMessage(conn net.Conn) {
//	reader := bufio.NewReader(conn)
//	for {
//		msg, err := reader.ReadString('\n')
//		if err != nil {
//			log.Fatal("程序错误! 服务器连接断开:", err)
//		}
//		fmt.Println(strings.TrimSpace(msg))
//	}
//}
