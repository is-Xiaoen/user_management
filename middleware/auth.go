package middleware

import (
	"net/http"

	"user-management-system/app"
	"user-management-system/errors"
	"user-management-system/session"
)

// Middleware 中间件集合
type Middleware struct {
	Auth *AuthMiddleware
}

// NewMiddleware 创建中间件集合
func NewMiddleware(application *app.App) *Middleware {
	sessionHelper := session.NewHelper(application.SessionManager, application.UserRepository)

	return &Middleware{
		Auth: &AuthMiddleware{
			app:           application,
			sessionHelper: sessionHelper,
		},
	}
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	app           *app.App
	sessionHelper *session.Helper
}

// RequireAuth 要求用户已登录
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 使用会话管理器检查用户是否已登录
		_, err := m.sessionHelper.RequireLogin(r)
		if err != nil {
			// 未登录，重定向到登录页面
			http.Redirect(w, r, "/login", http.StatusSeeOther)
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
		user, err := m.sessionHelper.GetCurrentUser(r)
		if err != nil {
			// 如果获取用户信息失败，重定向到登录页面
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// 使用服务层检查用户是否为管理员
		if !m.app.UserService.IsAdmin(user) {
			// 如果不是管理员，返回权限错误
			errors.HandleError(w, r, errors.NewForbiddenError("需要管理员权限"))
			return
		}

		// 用户是管理员，继续执行后续处理程序
		next.ServeHTTP(w, r)
	})
}
