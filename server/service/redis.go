package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
)

// AddUserActivity 增加用户活跃度
func AddUserActivity(client *Client) error {
	err := RDB.ZIncrBy(context.Background(), "chat:activity", 1, client.Username).Err()
	if err != nil {
		return fmt.Errorf("更新用户 %s 活跃度失败: %v", client.Username, err)
	}
	return nil
}

// GetRankData 获取用户活跃度排名数据
func GetRankData(client *Client) ([]string, error) {
	result, err := RDB.ZRevRange(context.Background(), "chat:activity", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("用户 %s 获取活跃度排名失败: %v", client.Username, err)
	}
	return result, nil
}
func GetScore(username string) float64 {
	score, err := RDB.ZScore(context.Background(), "chat:activity", username).Result()
	if err != nil {
		log.Printf("获取用户 %s score失败: %v", username, err)
		return 0
	}
	return score
}
func InitConsumerGroup() {
	err := RDB.XGroupCreateMkStream(context.Background(), streamNameChat, groupName, "$").Err()
	if err != nil {
		if strings.Contains(err.Error(), "BUSYGROUP") {
			return
		}
		log.Printf("创建消费者组失败: %v", err)
	}
}

// AddToStream 向Stream中发送消息
func (msg *ChatMessage) AddToStream(streamName string) error {
	msgData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json序列化失败: %w", err)
	}
	err = RDB.XAdd(context.Background(), &redis.XAddArgs{
		Stream: streamName,
		MaxLen: 10000, // 限制10000条消息记录
		Approx: true,  // 近似删除
		Values: map[string]interface{}{
			"message": msgData,
		},
	}).Err()
	if err != nil {
		return fmt.Errorf("消息写入stream失败: %w", err)
	}
	return nil
}
func GetHistoryMessage() ([]ChatMessage, error) {
	// 展示10条历史消息
	result, err := RDB.XRevRangeN(context.Background(), streamNameChat, "+", "-", 10).Result()
	if err != nil {
		return nil, fmt.Errorf("读取历史消息失败: %w", err)
	}
	var historyMessage []ChatMessage
	for _, message := range result {
		// 从stream获取消息
		messageData := message.Values["message"].(string)
		// 反序列化为chatMessage结构体类型
		var chatMsg ChatMessage
		err := json.Unmarshal([]byte(messageData), &chatMsg)
		if err != nil {
			log.Printf("消息反序列化失败: %v", err)
			continue
		}
		historyMessage = append(historyMessage, chatMsg)
	}
	return historyMessage, nil
}
