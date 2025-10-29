package service

import (
	"log"
	"time"
)

func HeartBeatCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		lock.Lock()
		for conn, client := range users {
			// client内部锁保证读取时间正确
			client.mu.Lock()
			duration := now.Sub(client.LastActiveTime)
			client.mu.Unlock()
			if duration > 2*time.Minute {
				conn.Write([]byte("[系统] 超时: 您长时间未任何操作! 连接已自动断开!\n"))
				client.Conn.Close()
				delete(users, client.Conn)
				log.Printf("[系统] 用户 %s 长时间未活跃, 自动断开", client.Username)
			} else if duration > 1*time.Minute {
				conn.Write([]byte("[系统] 注意: 请保持活跃状态!\n"))
			}
		}
		lock.Unlock()
	}
}
