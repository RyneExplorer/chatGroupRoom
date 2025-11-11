package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"runtime/debug"
)

// BroadcastFromStreamLoop 从Stream中读取消息并发送给客户端
func BroadcastFromStreamLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("stream群聊消息出现panic... %v\n%s", err, debug.Stack())
		}
	}()
	for {
		result, err := RDB.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName1,
			Streams:  []string{streamNameChat, ">"},
			Count:    10,
			Block:    0,
		}).Result()
		if err != nil {
			if err == redis.Nil { // 查询到不存在的key或超时未收到消息
				continue
			}
			log.Printf("读取stream消息失败: %v", err)
			return
		}

		for _, stream := range result {
			for _, message := range stream.Messages {
				var chatMsg ChatMessage
				err := json.Unmarshal([]byte(message.Values["message"].(string)), &chatMsg)
				if err != nil {
					log.Printf("消息反序列化失败: %v", err)
					continue
				}

				switch chatMsg.Type {
				case MessageTypePublic:
					broadcastToAll(fmt.Sprintf("[群聊] %s: %s", chatMsg.Username, chatMsg.Message))
				}
				RDB.XAck(context.Background(), streamNameChat, groupName, message.ID)
			}
		}
	}
}
