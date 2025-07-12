package controllers

import (
	"html/template"
	"strings"

	"user-management-system/app"
)

// Controllers 控制器集合
type Controllers struct {
	Auth *AuthController
	User *UserController
}

// NewControllers 创建控制器集合
// 注意：不再在这里初始化服务，而是让每个控制器自己管理
func NewControllers(application *app.App) *Controllers {
	return &Controllers{
		Auth: NewAuthController(application),
		User: NewUserController(application),
	}
}

// funcMap 定义模板函数
var funcMap = template.FuncMap{
	"upper": strings.ToUpper,
}
