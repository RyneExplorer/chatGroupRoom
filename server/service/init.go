package service

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	DB  *sql.DB
	RDB *redis.Client
)

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/chat")
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   2,
	})

	// 优先重置数据库用户状态
	err = ResetAllUserStatus(LEAVE_STATUS)
	if err != nil {
		log.Printf("重置用户状态失败: %v", err)
		return
	}

	// 创建消费者组
	InitConsumerGroup()
}
