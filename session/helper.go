package session

import (
	"fmt"
	"net/http"

	"user-management-system/errors"
	"user-management-system/models"
	"user-management-system/repository/interfaces"
)

// Helper 会话辅助器，封装常用操作
type Helper struct {
	manager        *Manager
	userRepository interfaces.UserRepository
}

// NewHelper 创建会话辅助器
func NewHelper(manager *Manager, userRepo interfaces.UserRepository) *Helper {
	return &Helper{
		manager:        manager,
		userRepository: userRepo,
	}
}

// GetCurrentUser 从请求中获取当前登录用户
func (h *Helper) GetCurrentUser(r *http.Request) (*models.User, error) {
	// 获取会话
	session, err := h.manager.GetSession(r)
	if err != nil {
		return nil, errors.NewUnauthorizedError("会话无效或已过期")
	}

	// 从会话中获取用户ID
	userID := session.UserID

	// 根据用户ID获取用户信息
	user, err := h.userRepository.GetByID(userID)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("获取用户信息失败: %w", err))
	}

	if user == nil {
		return nil, errors.NewNotFoundError("用户")
	}

	return user, nil
}

// Login 处理用户登录，创建会话
func (h *Helper) Login(w http.ResponseWriter, r *http.Request, userID int, remember bool) error {
	// 重要：先销毁旧会话，防止会话固定攻击
	h.manager.DestroySession(w, r)
	// 创建新会话
	_, err := h.manager.CreateSession(w, userID, remember)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("创建会话失败: %w", err))
	}
	return nil
}

// Logout 处理用户登出，销毁会话
func (h *Helper) Logout(w http.ResponseWriter, r *http.Request) {
	h.manager.DestroySession(w, r)
}

// RequireLogin 检查用户是否已登录
func (h *Helper) RequireLogin(r *http.Request) (*Session, error) {
	session, err := h.manager.GetSession(r)
	if err != nil {
		return nil, errors.NewUnauthorizedError("请先登录")
	}
	return session, nil
}

// GetCSRFTokenForTemplate 为模板获取CSRF令牌
func (h *Helper) GetCSRFTokenForTemplate(r *http.Request) (string, error) {
	session, err := h.manager.GetSession(r)
	if err != nil {
		return "", errors.NewUnauthorizedError("会话无效")
	}

	token, err := GetCSRFToken(session)
	if err != nil {
		return "", errors.NewInternalError(fmt.Errorf("获取CSRF令牌失败: %w", err))
	}

	return token, nil
}
