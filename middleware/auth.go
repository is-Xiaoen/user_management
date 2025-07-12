package middleware

import (
	"net/http"
	"sync"

	"user-management-system/app"
	"user-management-system/errors"
	"user-management-system/repository/mysql"
	"user-management-system/session"
)

// Middleware 中间件集合
type Middleware struct {
	Auth *AuthMiddleware
}

// NewMiddleware 创建中间件集合
func NewMiddleware(application *app.App) *Middleware {
	return &Middleware{
		Auth: NewAuthMiddleware(application),
	}
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	app           *app.App
	sessionHelper *session.Helper
	once          sync.Once
	mu            sync.RWMutex
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(application *app.App) *AuthMiddleware {
	return &AuthMiddleware{
		app: application,
	}
}

// getSessionHelper 延迟初始化会话助手
func (m *AuthMiddleware) getSessionHelper() *session.Helper {
	m.once.Do(func() {
		// 创建用户仓库（只用于会话助手）
		userRepo := mysql.NewUserRepository(m.app.GetDB())

		// 创建会话助手
		m.sessionHelper = session.NewHelper(m.app.GetSessionManager(), userRepo)
	})

	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessionHelper
}

// RequireAuth 要求用户已登录
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 使用会话管理器检查用户是否已登录
		sessionHelper := m.getSessionHelper()
		_, err := sessionHelper.RequireLogin(r)
		if err != nil {
			// 未登录，重定向到登录页面
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// RequirePostAuth 要求用户已登录
func (m *AuthMiddleware) RequirePostAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 使用会话管理器检查用户是否已登录
		sessionHelper := m.getSessionHelper()
		_, err := sessionHelper.RequireLogin(r)
		if err != nil {
			// 已经登录，重定向到管理页面
			http.Redirect(w, r, "/users", http.StatusSeeOther)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin 要求管理员权限
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取当前用户
		sessionHelper := m.getSessionHelper()
		user, err := sessionHelper.GetCurrentUser(r)
		if err != nil {
			// 如果获取用户信息失败，重定向到登录页面
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// 检查用户是否为管理员
		if !user.IsAdmin() {
			// 如果不是管理员，返回权限错误
			errors.HandleError(w, r, errors.NewForbiddenError("需要管理员权限"))
			return
		}

		// 用户是管理员，继续执行后续处理程序
		next.ServeHTTP(w, r)
	})
}
