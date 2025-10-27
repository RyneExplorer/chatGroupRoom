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
}
