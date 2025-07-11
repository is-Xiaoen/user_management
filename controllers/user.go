package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"user-management-system/app"
	"user-management-system/errors"
	"user-management-system/logger"
	"user-management-system/models"
	"user-management-system/session"
)

// UserController 用户管理控制器
type UserController struct {
	app           *app.App
	sessionHelper *session.Helper
}

// NewUserController 创建用户控制器
func NewUserController(application *app.App, sessionHelper *session.Helper) *UserController {
	return &UserController{
		app:           application,
		sessionHelper: sessionHelper,
	}
}

// RenderHomePage 渲染首页
func (c *UserController) RenderHomePage(w http.ResponseWriter, r *http.Request) {
	// 获取当前用户（可选）
	currentUser, _ := c.sessionHelper.GetCurrentUser(r)

	data := struct {
		CurrentUser *models.User
	}{
		CurrentUser: currentUser,
	}

	tmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("views/layout.html", "views/index.html")
	if err != nil {
		log.Printf("模板解析错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("模板执行错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
	}
}

// RenderUsersPage 渲染用户列表页面
func (c *UserController) RenderUsersPage(w http.ResponseWriter, r *http.Request) {
	// 获取当前用户
	currentUser, err := c.sessionHelper.GetCurrentUser(r)
	if err != nil {
		errors.HandleError(w, r, errors.NewUnauthorizedError(""))
		return
	}

	// 记录查看用户列表操作
	logger.UserAction(currentUser.Username, "查看用户列表", "", true)

	// 获取所有用户
	users, err := c.app.UserService.GetAllUsers()
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	// 获取CSRF令牌
	csrfToken, err := c.sessionHelper.GetCSRFTokenForTemplate(r)
	if err != nil {
		log.Printf("获取CSRF令牌失败: %v", err)
		csrfToken = "" // 继续处理，但不使用CSRF保护
	}

	data := struct {
		CurrentUser *models.User
		Users       []*models.User
		CSRFToken   string
	}{
		CurrentUser: currentUser,
		Users:       users,
		CSRFToken:   csrfToken,
	}

	tmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("views/layout.html", "views/users.html")
	if err != nil {
		log.Printf("模板解析错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("模板执行错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
	}
}

// HandleDeleteUser  处理删除用户请求
func (c *UserController) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.HandleError(w, r, errors.NewAppError(
			errors.ValidationError,
			"方法不允许",
			nil,
		))
		return
	}

	//解析表单
	if err := r.ParseForm(); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "无法解析表单"))
		return
	}

	//获取用户ID
	userIDStr := r.FormValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "无效的用户ID"))
		return
	}

	//获取当前用户
	currentUser, err := c.sessionHelper.GetCurrentUser(r)
	//currentUser, err := session.GetCurrentUser(r)
	if err != nil {
		errors.HandleError(w, r, errors.NewUnauthorizedError(""))
		return
	}

	//防止删除自己
	if userID == currentUser.ID {
		errors.HandleError(w, r, errors.NewForbiddenError("不能删除自己"))
		return
	}

	// 获取要删除的用户信息（用于日志记录）
	targetUser, _ := c.app.UserService.GetUserByID(userID)
	targetUsername := targetUser.Username

	//删除用户
	if err := c.app.UserService.DeleteUser(userID); err != nil {
		// 记录删除失败
		logger.UserActionWithError(currentUser.Username, "删除用户",
			fmt.Sprintf("目标用户: %s (ID: %d)", targetUsername, userID), err)
		errors.HandleError(w, r, err)
		return
	}

	// 记录删除成功
	logger.UserAction(currentUser.Username, "删除用户",
		fmt.Sprintf("目标用户: %s (ID: %d)", targetUsername, userID), true)
	//重新定向到用户列表
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// HandleUpdateUser 处理更新用户请求
func (c *UserController) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.HandleError(w, r, errors.NewAppError(
			errors.ValidationError,
			"方法不允许",
			nil,
		))
		return
	}

	//解析表单
	if err := r.ParseForm(); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "无法解析表单"))
		return
	}
	//获取表单数据
	userIDStr := r.FormValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "无效的用户ID"))
		return
	}
	email := r.FormValue("email")
	role := r.FormValue("role")

	//获取当前用户
	currentUser, err := c.sessionHelper.GetCurrentUser(r)
	if err != nil {
		errors.HandleError(w, r, errors.NewUnauthorizedError(""))
		return
	}

	//获取更新的用户信息 (记录日志)
	targetUser, _ := c.app.UserService.GetUserByID(userID)
	targetUsername := targetUser.Username

	//更新用户
	if err := c.app.UserService.UpdateUser(userID, email, role); err != nil {
		// 记录更新失败
		logger.UserActionWithError(currentUser.Username, "更新用户",
			fmt.Sprintf("目标用户: %s (ID: %d)", targetUsername, userID), err)
		errors.HandleError(w, r, err)
		return
	}

	// 记录更新成功
	logger.UserAction(currentUser.Username, "更新用户",
		fmt.Sprintf("目标用户: %s (ID: %d), 邮箱: %s, 角色: %s",
			targetUsername, userID, email, role), true)

	// 重定向到用户列表
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
