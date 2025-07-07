package middleware

import (
	"user-management-system/services"
)

// 全局服务实例
var userService services.UserService

// InitMiddleware 初始化中间件
func InitMiddleware(service *services.Service) {
	userService = service.UserService
}
