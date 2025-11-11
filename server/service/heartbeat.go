package service

import (
	"log"
	"runtime/debug"
	"time"
)

func HeartBeatCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("心跳监测协程出现panic... %v\n%s", err, debug.Stack())
		}
	}()
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		Lock.Lock()
		for conn, client := range Users {
			// client内部锁保证读取时间正确
			client.Mu.Lock()
			duration := now.Sub(client.LastActiveTime)
			client.Mu.Unlock()
			if duration > 20*time.Minute {
				//conn.Write([]byte("[系统] 超时: 你长时间未任何操作! 连接已自动断开!\n"))
				warning := "[系统] 超时: 你长时间未任何操作! 连接已自动断开!"
				err := writeMessage(conn, warning)
				if err != nil {
					log.Printf("提醒用户 %s 失败: %v", client.Username, err)
				}
				delete(Users, client.Conn)
				client.Conn.Close()
				log.Printf("[系统] 用户 %s 长时间未活跃, 自动断开", client.Username)
			} else if duration > 15*time.Minute {
				//conn.Write([]byte("[系统] 注意: 请保持活跃状态!\n"))
				warning := "[系统] 注意: 请保持活跃状态!"
				err := writeMessage(conn, warning)
				if err != nil {
					log.Printf("提醒用户 %s 失败: %v", client.Username, err)
				}
			}
		}
		Lock.Unlock()
	}
}
