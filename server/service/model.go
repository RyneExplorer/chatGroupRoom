package service

import (
	"errors"
	"net"
	"sync"
	"time"
)

const (
	LEAVE_STATUS  = 0
	ONLINE_STATUS = 1
)

type User struct {
	ID         int       `db:"id" json:"ID"`
	Username   string    `db:"username" json:"username"`
	Password   string    `db:"password" json:"-"`
	Status     int       `db:"status" json:"status"`
	CreateTime time.Time `db:"create_time" json:"createTime"`
	LastTime   time.Time `db:"last_time" json:"lastTime"`
}
type Client struct {
	Username       string
	Conn           net.Conn
	LastActiveTime time.Time
	mu             sync.Mutex
}

var (
	users       = make(map[net.Conn]*Client)
	msgChan     = make(chan string, 10)
	onlineChan  = make(chan string, 5)
	leaveChan   = make(chan string, 5)
	privateChan = make(chan string, 5)
	listChan    = make(chan net.Conn, 5)
	lock        = sync.Mutex{} // 全局锁
)

var (
	ErrUsernameExists     = errors.New("用户名已存在")
	ErrUsernameNotExists  = errors.New("用户名已存在")
	ErrInvalidInput       = errors.New("用户名或密码不能为空")
	ErrUsernamePassData   = errors.New("密码错误")
	ErrUserIsOnline       = errors.New("用户已经上线")
	ErrNotSupported       = errors.New("不支持该操作")
	ErrDataFormat         = errors.New("数据格式错误")
	ErrUpdateStatusFailed = errors.New("更新用户状态失败")
)
