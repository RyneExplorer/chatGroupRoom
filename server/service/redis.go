package service

import (
	"context"
	"fmt"
	"log"
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
