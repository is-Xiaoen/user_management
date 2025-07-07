package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"user-management-system/errors"
	"user-management-system/logger"
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

// RenderRegisterPage 渲染注册页面
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
	user, err := userService.AuthenticateUser(username, password)
	if err != nil {
		// 记录登录失败    UserAction是记录用户操作日志
		//fmt.Println("controllers\auth.go 的 119行")
		//fmt.Println(err)
		//fmt.Println(1111)
		logger.UserAction(username, "登录", "IP: "+r.RemoteAddr, false)

		//渲染登录页面并显示错误
		appErr, _ := errors.IsAppError(err)

		data := struct {
			CurrentUser *models.User
			Error       string
		}{
			CurrentUser: nil,
			Error:       appErr.Message,
		}

		tmpl, parseErr := template.New("layout.html").Funcs(funcMap).ParseFiles("views/layout.html", "views/login.html")
		if parseErr != nil {
			errors.HandleError(w, r, errors.NewInternalError(parseErr))
			return
		}

		if execErr := tmpl.ExecuteTemplate(w, "layout", data); execErr != nil {
			errors.HandleError(w, r, errors.NewInternalError(execErr))
		}
		return
	}

	// 使用会话管理创建会话
	if err := session.Login(w, user.ID, remember); err != nil {
		errors.HandleError(w, r, errors.NewInternalError(err))
		return
	}

	//记录登录成功
	logger.UserAction(username, "登录", "IP: "+r.RemoteAddr, true)

	//登录成功后,重定向到用户列表页面
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// HandleRegister 处理用户注册
func HandleRegister(w http.ResponseWriter, r *http.Request) {
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
	email := r.FormValue("email")

	//使用服务层注册用户
	err = userService.RegisterUser(username, password, email)
	if err != nil {
		//记录注册失败
		logger.UserAction(username, "注册", "邮箱: "+email+", IP: "+r.RemoteAddr, false)

		//渲染注册页面并希纳是错误信息
		appErr, _ := errors.IsAppError(err)
		data := struct {
			CurrentUser *models.User
			Error       string
		}{
			CurrentUser: nil,
			Error:       appErr.Message,
		}

		tmpl, parseErr := template.New("layout.html").Funcs(funcMap).ParseFiles("views/layout.html", "views/register.html")
		if parseErr != nil {
			errors.HandleError(w, r, errors.NewInternalError(parseErr))
			return
		}

		if execErr := tmpl.ExecuteTemplate(w, "layout", data); execErr != nil {
			errors.HandleError(w, r, errors.NewInternalError(execErr))
		}
		return
	}
	//记录注册成功
	logger.UserAction(username, "注册", "邮箱: "+email+", IP: "+r.RemoteAddr, true)

	//注册成功后 重新定向到登录页面
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// HandleLogout  处理用户登出
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	//获取当前用户信息  ,用于记录日志
	currentUser, _ := session.GetCurrentUser(r)

	//使用会话管理器销毁会话
	session.Logout(w, r)

	//记录登出
	if currentUser != nil {
		logger.UserAction(currentUser.Username, "登出", "IP: "+r.RemoteAddr, true)
	}

	//清除会话后, 重新定向到登录页面
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
