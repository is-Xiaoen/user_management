package controllers

import (
	"html/template"
	"log"
	"net/http"
	"sync"

	"user-management-system/app"
	"user-management-system/errors"
	"user-management-system/logger"
	"user-management-system/models"
	_ "user-management-system/repository/interfaces"
	"user-management-system/repository/mysql"
	"user-management-system/services"
	"user-management-system/session"
)

// AuthController 认证控制器
type AuthController struct {
	app           *app.App
	sessionHelper *session.Helper
	userService   services.UserService
	once          sync.Once    // 确保服务只初始化一次
	mu            sync.RWMutex // 保护并发访问
}

// NewAuthController 创建认证控制器
func NewAuthController(application *app.App) *AuthController {
	return &AuthController{
		app: application,
		// 注意：这里不再立即初始化服务
	}
}

// getUserService 延迟初始化用户服务
func (c *AuthController) getUserService() services.UserService {
	c.once.Do(func() {
		// 创建用户仓库
		userRepo := mysql.NewUserRepository(c.app.GetDB())

		// 创建用户服务
		c.userService = services.NewUserService(userRepo)

		// 创建会话助手
		c.sessionHelper = session.NewHelper(c.app.GetSessionManager(), userRepo)

		logger.Info("AuthController: 用户服务已初始化")
	})

	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userService
}

// getSessionHelper 获取会话助手
func (c *AuthController) getSessionHelper() *session.Helper {
	// 确保服务已初始化
	c.getUserService()

	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sessionHelper
}

// RenderLoginPage 渲染登录页面
func (c *AuthController) RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	// 使用延迟初始化的会话助手
	sessionHelper := c.getSessionHelper()

	// 尝试获取会话，如果会话存在，则用户已登录
	_, err := sessionHelper.RequireLogin(r)
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

// HandleLogin 处理用户登录
func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法是否为 POST
	if r.Method != http.MethodPost {
		errors.HandleError(w, r, errors.NewAppError(
			errors.ValidationError,
			"方法不允许",
			nil,
		))
		return
	}

	_, err := c.sessionHelper.RequireLogin(r)
	if err == nil {
		// 用户已登录，重定向到用户列表页面
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	// 解析 HTTP 请求中的表单数据
	err = r.ParseForm()
	if err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "无法解析表单"))
		return
	}

	// 从表单中获取数据
	username := r.FormValue("username")
	password := r.FormValue("password")
	remember := r.FormValue("remember") == "on"

	// 使用延迟初始化的服务层验证用户
	userService := c.getUserService()
	user, err := userService.AuthenticateUser(username, password)
	if err != nil {
		// 记录登录失败
		logger.UserAction(username, "登录", "IP: "+r.RemoteAddr, false)

		// 渲染登录页面并显示错误信息
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

	// 使用会话管理器创建会话
	sessionHelper := c.getSessionHelper()
	if err := sessionHelper.Login(w, r, user.ID, remember); err != nil {
		errors.HandleError(w, r, errors.NewInternalError(err))
		return
	}

	// 记录登录成功
	logger.UserAction(user.Username, "登录", "IP: "+r.RemoteAddr, true)

	// 登录成功后，重定向到用户列表页面
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// HandleRegister 处理用户注册
func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法是否为 POST
	if r.Method != http.MethodPost {
		errors.HandleError(w, r, errors.NewAppError(
			errors.ValidationError,
			"方法不允许",
			nil,
		))
		return
	}

	// 解析 HTTP 请求中的表单数据
	err := r.ParseForm()
	if err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "无法解析表单"))
		return
	}

	// 从表单中获取数据
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	// 使用延迟初始化的服务层注册用户
	userService := c.getUserService()
	err = userService.RegisterUser(username, password, email)
	if err != nil {
		// 记录注册失败
		logger.UserAction(username, "注册", "邮箱: "+email+", IP: "+r.RemoteAddr, false)

		// 渲染注册页面并显示错误信息
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

	// 记录注册成功
	logger.UserAction(username, "注册", "邮箱: "+email+", IP: "+r.RemoteAddr, true)

	// 注册成功后，重定向到登录页面
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// HandleLogout 处理用户登出
func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionHelper := c.getSessionHelper()

	// 获取当前用户信息（用于记录日志）
	currentUser, _ := sessionHelper.GetCurrentUser(r)

	// 使用会话管理器销毁会话
	sessionHelper.Logout(w, r)

	// 记录登出
	if currentUser != nil {
		logger.UserAction(currentUser.Username, "登出", "IP: "+r.RemoteAddr, true)
	}

	// 清除会话后，重定向到登录页面
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// RenderRegisterPage 渲染注册页面（保持原有实现）
func (c *AuthController) RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	sessionHelper := c.getSessionHelper()

	// 尝试获取会话，如果会话存在，则用户已登录
	_, err := sessionHelper.RequireLogin(r)
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
