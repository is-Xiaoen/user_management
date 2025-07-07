package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-management-system/logger"

	"user-management-system/errors"
	"user-management-system/models"
	"user-management-system/session"
)

// funcMap 定义模板函数
var funcMap = template.FuncMap{
	"upper": strings.ToUpper,
}

// RenderHomePage 渲染主页
func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	//获取当前用户
	currentUser, _ := session.GetCurrentUser(r)

	data := struct {
		CurrentUser *models.User
	}{
		CurrentUser: currentUser,
	}

	tmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("views/layout", "views/index.html")

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
func RenderUsersPage(w http.ResponseWriter, r *http.Request) {
	//获取当前用户
	currentUser, err := session.GetCurrentUser(r)
	if err != nil {
		errors.HandleError(w, r, errors.NewUnauthorizedError(""))
		return
	}

	//记录查看用户列表操作
	logger.UserAction(currentUser.Username, "查看用户列表", "", true)

	//获取所有用户
	users, err := userService.GetAllUsers()
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	// 获取CSRF令牌
	csrfToken, err := session.GetCSRFTokenForTemplate(r)
	if err != nil {
		log.Printf("获取CSRF令牌失败: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(fmt.Errorf("获取CSRF令牌失败: %w", err)))
		return
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
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
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
	currentUser, err := session.GetCurrentUser(r)
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
	targetUser, _ := userService.GetUserByID(userID)
	targetUsername := targetUser.Username

	//删除用户
	if err := userService.DeleteUser(userID); err != nil {
		// 记录删除失败
		logger.UserAction(currentUser.Username, "删除用户", fmt.Sprintf("目标用户: %s (ID: %d), 失败原因: %v", targetUsername, userID, err), false)
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
func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
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
	currentUser, err := session.GetCurrentUser(r)
	if err != nil {
		errors.HandleError(w, r, errors.NewUnauthorizedError(""))
		return
	}

	//获取更新的用户信息 (记录日志)
	targetUser, _ := userService.GetUserByID(userID)
	targetUsername := targetUser.Username

	//更新用户
	if err := userService.UpdateUser(userID, email, role); err != nil {
		// 记录更新失败
		logger.UserAction(currentUser.Username, "更新用户",
			fmt.Sprintf("目标用户: %s (ID: %d), 邮箱: %s, 角色: %s, 失败原因: %v",
				targetUsername, userID, email, role, err), false)
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

// HandleAPIUsers 处理API用户列表请求
func HandleAPIUsers(w http.ResponseWriter, r *http.Request) {
	// 获取所有用户
	users, err := userService.GetAllUsers()
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	// 转换为安全的API响应格式
	type APIUser struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	apiUsers := make([]APIUser, len(users))
	for i, user := range users {
		apiUsers[i] = APIUser{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		}
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 编码并发送响应
	if err := json.NewEncoder(w).Encode(apiUsers); err != nil {
		log.Printf("JSON编码错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
	}
}

// HandleAPIUserStats 处理API用户统计请求
func HandleAPIUserStats(w http.ResponseWriter, r *http.Request) {
	// 获取用户统计信息
	stats, err := userService.GetUserStats()
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 编码并发送响应
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("JSON编码错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
	}
}
