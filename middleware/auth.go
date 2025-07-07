package middleware

import (
	"net/http"
	"user-management-system/errors"
	"user-management-system/services"
	"user-management-system/session"
)

// 全局服务实例
var userService services.UserService

// InitMiddleware 初始化中间件
func InitMiddleware(service *services.Service) {
	userService = service.UserService
}

// AuthMiddleware 验证用户是否已经登录
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//使用会话管理器检查用户是否已经登录
		_, err := session.RequireLogin(r)
		if err != nil {
			//未登录，重定向到登录页面
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		//继续处理登录请求
		next.ServeHTTP(w, r)
	})
}

// AdminMiddleware 验证用户是否是管理员
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//获取当前用户
		user, err := session.GetCurrentUser(r)
		if err != nil {
			// 如果获取用户信息失败，重定向到登录页面
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// 使用服务层见擦汗用户是否为管理员
		if !user.IsAdmin() {
			// 如果不是管理员，返回权限错误
			errors.HandleError(w, r, errors.NewForbiddenError("需要管理员权限"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
