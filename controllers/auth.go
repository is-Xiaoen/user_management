package controllers

import (
	"html/template"
	"log"
	"net/http"
	"user-management-system/errors"
	"user-management-system/models"
	"user-management-system/services"
	"user-management-system/session"
)

// 全局服务实例
var userService services.UserService

// InitControllers 初始化控制器
func InitControllers(service *services.Service) {
	userService = service.UserService
}

// RenderLoginPage 渲染登陆页面
func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	// 尝试获取对话 , 如果对话存在,则用户已登录
	_, err := session.RequireLogin(r)
	if err == nil {
		// 用户已登录，重定向到用户列表页面
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	// 准备传递给模板的数据
	data := struct {
		CurrentUser *models.User
		Error       string
	}{
		CurrentUser: nil,
		Error:       "",
	}

	// 解析模板文件
	tmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("views/layout.html", "views/login.html")
	if err != nil {
		log.Printf("模板解析错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
		return
	}
	// 执行模板渲染
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("模板执行错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
	}
}

// RenderRegisterPage渲染注册页面
func RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	// 尝试获取对话 , 如果对话存在,则用户已登录
	_, err := session.RequireLogin(r)
	if err == nil {
		// 用户已登录，重定向到用户列表页面
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	// 准备传递给模板的数据
	data := struct {
		CurrentUser *models.User
		Error       string
	}{
		CurrentUser: nil,
		Error:       "",
	}

	// 解析注册页面所需的模板文件
	tmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("views/layout.html", "views/register.html")
	if err != nil {
		log.Printf("模板解析错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
		return
	}

	// 执行模板渲染
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("模板执行错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
	}
}

// HandleLogin 处理用户登录
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	//检查请求方法是否为POST
	if r.Method != http.MethodPost {
		errors.HandleError(w, r, errors.NewAppError(
			errors.ValidationError,
			"方法不允许",
			nil,
		))
		return
	}

	//解析 HTTP 请求中的表单数据
	err := r.ParseForm()
	if err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "无法解析表单"))
		return
	}
	//从表单中获取数据
	username := r.FormValue("username")
	password := r.FormValue("password")
	remember := r.FormValue("remember") == "on"

	//使用服务层验证用户
	user, err := userService.
}
