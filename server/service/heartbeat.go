package service

import (
	"log"
	"time"
)

func HeartBeatCheck() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		Lock.Lock()
		for conn, client := range Users {
			// client内部锁保证读取时间正确
			client.Mu.Lock()
			duration := now.Sub(client.LastActiveTime)
			client.Mu.Unlock()
			if duration > 10*time.Minute {
				conn.Write([]byte("[系统] 超时: 你长时间未任何操作! 连接已自动断开!\n"))
				client.Conn.Close()
				delete(Users, client.Conn)
				log.Printf("[系统] 用户 %s 长时间未活跃, 自动断开", client.Username)
			} else if duration > 15*time.Minute {
				conn.Write([]byte("[系统] 注意: 请保持活跃状态!\n"))
			}
		}
		Lock.Unlock()
	}
}
