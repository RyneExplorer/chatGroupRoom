package service

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

func writeMessage(conn io.Writer, msg string) error {
	msgBytes := []byte(msg)

	if len(msgBytes) > 65535 {
		return fmt.Errorf("消息过长: %d 字节", len(msgBytes))
	}

	// 1. 写入4字节长度（大端序）
	length := uint32(len(msgBytes))
	if err := binary.Write(conn, binary.BigEndian, length); err != nil {
		return fmt.Errorf("写入消息长度失败: %w", err)
	}

	// 2. 写入消息体
	_, err := conn.Write(msgBytes)
	if err != nil {
		return fmt.Errorf("写入消息内容失败: %w", err)
	}
	return nil
}

func readMessage(reader *bufio.Reader) (string, error) {
	// 先读取消息长度
	lengthBytes := make([]byte, 4)
	_, err := io.ReadFull(reader, lengthBytes)
	if err != nil {
		return "", fmt.Errorf("读取长度前缀失败: %w", err)
	}

	// 获取消息的长度
	msgLength := int(binary.BigEndian.Uint32(lengthBytes))
	if msgLength <= 0 {
		return "", fmt.Errorf("非法消息长度: %d", msgLength)
	}

	// 根据长度读取消息内容
	msgBytes := make([]byte, msgLength)
	_, err = io.ReadFull(reader, msgBytes)
	if err != nil {
		return "", fmt.Errorf("读取消息内容失败: %w", err)
	}
	return string(msgBytes), nil
}
