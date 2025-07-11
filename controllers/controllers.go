package controllers

import (
	"html/template"
	"strings"

	"user-management-system/app"
	"user-management-system/session"
)

// Controllers 控制器集合
type Controllers struct {
	Auth *AuthController
	User *UserController
}

// NewControllers 创建控制器集合
func NewControllers(application *app.App) *Controllers {
	sessionHelper := session.NewHelper(application.SessionManager, application.UserRepository)

	return &Controllers{
		Auth: NewAuthController(application, sessionHelper),
		User: NewUserController(application, sessionHelper),
	}
}

// funcMap 定义模板函数
var funcMap = template.FuncMap{
	"upper": strings.ToUpper,
}
