package service

import (
	"database/sql"
	"errors"
	"fmt"
)

// Register 注册功能检查数据库
func Register(username, password string) error {
	if username == "" || password == "" {
		return ErrInvalidInput
	}
	// 用户名和密码校验
	if len(username) < 3 || len(username) > 12 {
		return errors.New("用户名长度在3到12位之间")
	}
	if len(password) < 6 || len(password) > 18 {
		return errors.New("密码长度在6到18位之间")
	}
	// 检查数据库用户名是否已存在
	var count int
	row := DB.QueryRow("select count(*) from user where username = ?;", username)
	err := row.Scan(&count)
	if err != nil {
		return fmt.Errorf("查询数据库失败: %w", err)
	}
	if count > 0 {
		return ErrUsernameExists
	}
	// 插入数据
	_, err = DB.Exec("insert into user (username, password, create_time) values (?, ?, NOW());", username, password)
	if err != nil {
		return fmt.Errorf("写入用户数据失败: %w", err)
	}
	err = UpdateStatus(username, ONLINE_STATUS)
	if err != nil {
		return ErrUpdateStatusFailed
	}
	// 注册成功
	return nil
}

// Login 登录功能检查数据库
func Login(username, password string) error {
	if username == "" || password == "" {
		return ErrInvalidInput
	}
	var dbPassword string
	var status int
	err := DB.QueryRow("select password, status from user where username = ?;", username).Scan(&dbPassword, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUsernameNotExists
		}
		return fmt.Errorf("查询数据库失败: %w", err)
	}
	if status == ONLINE_STATUS {
		return ErrUserIsOnline
	}
	if dbPassword != password {
		return ErrUsernamePassData
	}
	// 更新登录时间和用户状态
	_, err = DB.Exec("update user set last_login_time = NOW() where username = ?;", username)
	if err != nil {
		return fmt.Errorf("更新时间数据失败: %w", err)
	}
	err = UpdateStatus(username, ONLINE_STATUS)
	if err != nil {
		return ErrUpdateStatusFailed
	}
	return nil
}
func UpdateStatus(username string, status int) error {
	_, err := DB.Exec("update user set status = ? where username = ?;", status, username)
	if err != nil {
		return ErrUpdateStatusFailed
	}
	return nil
}
func ResetAllUserStatus(status int) error {
	_, err := DB.Exec("update user set status = ? where status = ?;", status, ONLINE_STATUS)
	if err != nil {
		return ErrUpdateStatusFailed
	}
	return nil
}
