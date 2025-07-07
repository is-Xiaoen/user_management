package session

import (
	"fmt"
	"net/http"
	"user-management-system/errors"
	"user-management-system/models"
	"user-management-system/repository/interfaces"
)

// 全局仓储实例
var userRepository interfaces.UserRepository

// InitSession 初始化会话模块
func InitSession(userRepo interfaces.UserRepository) {
	userRepository = userRepo
}

// GetCurrentUser 从请求中获取当前登录用户
func GetCurrentUser(r *http.Request) (*models.User, error) {
	//获取会话
	session, err := GlobalManager.GetSession(r)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("获取用户信息失败: %w", err))
	}

	//从会话中获取用户ID
	userID := session.UserID

	//根据用户ID获取用户信息
	user, err := userRepository.GetByID(userID)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("获取用户信息失败: %w", err))
	}

	if user == nil {
		return nil, errors.NewNotFoundError("用户")
	}

	return user, nil
}

// Login 处理用户登录, 创建会话
func Login(w http.ResponseWriter, userID int, remember bool) error {
	//创建会话
	_, err := GlobalManager.CreateSession(w, userID, remember)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("创建会话失败: %w", err))
	}
	return nil
}

// Logout 处理用户登出，销毁会话
func Logout(w http.ResponseWriter, r *http.Request) {
	GlobalManager.DestroySession(w, r)
}

// SetSessionData 在会话中存储数据
func SetSessionData(r *http.Request, key string, value interface{}) error {
	session, err := GlobalManager.GetSession(r)
	if err != nil {
		return errors.NewUnauthorizedError("会话无效")
	}

	session.Data[key] = value
	return nil
}

// GetSessionData 从会话中获取数据
func GetSessionData(r *http.Request, key string) (interface{}, error) {
	session, err := GlobalManager.GetSession(r)
	if err != nil {
		return nil, errors.NewUnauthorizedError("会话无效")
	}

	value, ok := session.Data[key]
	if !ok {
		return nil, nil
	}

	return value, nil
}

// RequireLogin 检查用户是否登录
func RequireLogin(r *http.Request) (*Session, error) {
	session, err := GlobalManager.GetSession(r)
	if err != nil {
		return nil, errors.NewUnauthorizedError("未登录")
	}
	return session, nil
}

// GetCSRFTokenForTemplate 为模板获取CSRF令牌
func GetCSRFTokenForTemplate(r *http.Request) (string, error) {
	session, err := GlobalManager.GetSession(r)
	if err != nil {
		return "", errors.NewUnauthorizedError("会话无效")
	}

	token, err := GetCSRFToken(session)
	if err != nil {
		return "", errors.NewInternalError(fmt.Errorf("获取CSRF令牌失败: %w", err))
	}

	return token, nil
}
