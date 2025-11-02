package service

import (
	"errors"
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
	err := IsUsernameExist(username)
	if err != nil {
		return err
	}
	// 插入用户数据
	err = InsertNewUser(username, password)
	if err != nil {
		return err
	}
	// 注册成功
	return nil
}

// Login 登录功能检查数据库
func Login(username, password string) error {
	if username == "" || password == "" {
		return ErrInvalidInput
	}
	// 检验用户能否登录
	err := CanUserLogin(username, password)
	if err != nil {
		return err
	}
	// 更新登录时间和用户状态
	err = UpdateUserLoginTimeAndStatus(username)
	if err != nil {
		return err
	}
	return nil
}
