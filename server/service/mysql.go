package service

import (
	"database/sql"
	"errors"
	"fmt"
)

func IsUsernameExist(username string) error {
	var count int
	row := DB.QueryRow("select count(*) from user where username = ?;", username)
	err := row.Scan(&count)
	if err != nil {
		return fmt.Errorf("查询数据库失败: %w", err)
	}
	if count > 0 {
		return ErrUsernameExists
	}
	return nil
}
func InsertNewUser(username, password string) error {
	_, err := DB.Exec("insert into user (username, password, create_time) values (?, ?, NOW());", username, password)
	if err != nil {
		return fmt.Errorf("写入用户数据失败: %w", err)
	}
	err = UpdateStatus(username, ONLINE_STATUS)
	if err != nil {
		return ErrUpdateStatusFailed
	}
	return nil
}

func CanUserLogin(username, password string) error {
	var dbPassword string
	var status int
	// 校验用户密码是否正确
	err := DB.QueryRow("select password, status from user where username = ?;", username).Scan(&dbPassword, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUsernameNotExists
		}
		return fmt.Errorf("查询数据库失败: %w", err)
	}
	if dbPassword != password {
		return ErrUsernamePassData
	}
	// 校验用户状态是否正常
	if status == ONLINE_STATUS {
		return ErrUserIsOnline
	}
	return nil
}
func UpdateUserLoginTimeAndStatus(username string) error {
	_, err := DB.Exec("update user set last_login_time = NOW() where username = ?;", username)
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
func getUserLevel(username string) int {
	var level int
	err := DB.QueryRow("select level from user where username = ?;", username).Scan(&level)
	if err != nil {
		return 1
	}
	return level
}
func UpdateUserLevel(username string, level int) error {
	_, err := DB.Exec("update user set level = ? where username = ?;", level, username)
	if err != nil {
		return fmt.Errorf("更新用户 %s 等级失败: %w", username, err)
	}
	return nil
}
