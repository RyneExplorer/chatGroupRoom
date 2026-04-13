package user

//在代码里增或修改一些bug

import (
	"encoding/binary"
	"fmt"
	"net"
)

func writeMessage(conn net.Conn, msg string) error {
	msgBytes := []byte(msg)

	// 写入4字节大端长度
	if err := binary.Write(conn, binary.BigEndian, uint32(len(msgBytes))); err != nil {
		return fmt.Errorf("写入消息长度失败: %w", err)
	}

	// 写入消息体
	n, err := conn.Write(msgBytes)
	if n == 0 {
		return fmt.Errorf("消息内容为空")
	}
	if err != nil {
		return fmt.Errorf("写入消息内容失败: %w", err)
	}
	return fmt.Errorf("无错误")
}
