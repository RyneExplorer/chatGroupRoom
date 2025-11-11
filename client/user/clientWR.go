package user

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func writeMessage(conn net.Conn, msg string) error {
	msgBytes := []byte(msg)

	// 写入4字节大端长度
	if err := binary.Write(conn, binary.BigEndian, uint32(len(msgBytes))); err != nil {
		return fmt.Errorf("写入消息长度失败: %w", err)
	}

	// 写入消息体
	_, err := conn.Write(msgBytes)
	if err != nil {
		return fmt.Errorf("写入消息内容失败: %w", err)
	}
	return nil
}
func readMessage(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	// 读4字节长度
	lengthBytes := make([]byte, 4)
	if _, err := io.ReadFull(reader, lengthBytes); err != nil {
		return "", fmt.Errorf("读取消息长度失败: %w", err)
	}

	msgLength := int(binary.BigEndian.Uint32(lengthBytes))
	if msgLength <= 0 {
		return "", fmt.Errorf("非法消息长度: %d", msgLength)
	}

	// 读取消息体
	msgBytes := make([]byte, msgLength)
	if _, err := io.ReadFull(reader, msgBytes); err != nil {
		return "", fmt.Errorf("读取消息内容失败: %w", err)
	}

	return string(msgBytes), nil
}
