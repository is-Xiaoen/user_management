package router

import (
	"fmt"
	"net/http"

	"user-management-system/app"
	"user-management-system/controllers"
	"user-management-system/middleware"
	"user-management-system/session"
)

// Router 路由器
type Router struct {
	mux         *http.ServeMux
	app         *app.App
	controllers *controllers.Controllers
	middleware  *middleware.Middleware
}

// NewRouter 创建路由器
func NewRouter(application *app.App) *Router {
	return &Router{
		mux:         http.NewServeMux(),
		app:         application,
		controllers: controllers.NewControllers(application),
		middleware:  middleware.NewMiddleware(application),
	}
}

// Setup 设置路由
func (r *Router) Setup() http.Handler {
	// 创建CSRF中间件
	csrfMiddleware := session.NewCSRFMiddleware(r.app.SessionManager)

	// 静态文件
	fs := http.FileServer(http.Dir("static"))
	r.mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// 主页
	r.mux.HandleFunc("/", r.handleHome)

	// 认证相关
	r.mux.HandleFunc("/login", r.handleLogin)
	r.mux.HandleFunc("/register", r.handleRegister)
	r.mux.HandleFunc("/logout", r.controllers.Auth.HandleLogout)

	// 用户管理（需要认证）
	r.mux.Handle("/users", r.middleware.Auth.RequireAuth(
		http.HandlerFunc(r.controllers.User.RenderUsersPage),
	))

	// 用户删除（需要管理员权限 + CSRF保护）
	r.mux.Handle("/users/delete", r.middleware.Auth.RequireAdmin(
		csrfMiddleware(http.HandlerFunc(r.controllers.User.HandleDeleteUser)),
	))

	// 用户更新（需要管理员权限 + CSRF保护）
	r.mux.Handle("/users/update", r.middleware.Auth.RequireAdmin(
		csrfMiddleware(http.HandlerFunc(r.controllers.User.HandleUpdateUser)),
	))

	// API路由
	//r.mux.HandleFunc("/api/users", r.controllers.User.HandleAPIUsers)
	//r.mux.HandleFunc("/api/users/stats", r.controllers.User.HandleAPIUserStats)

	// 健康检查
	r.mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	return r.mux
}

func (r *Router) handleHome(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	r.controllers.User.RenderHomePage(w, req)
}

func (r *Router) handleLogin(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.controllers.Auth.HandleLogin(w, req)
	} else {
		r.controllers.Auth.RenderLoginPage(w, req)
	}
}

func (r *Router) handleRegister(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.controllers.Auth.HandleRegister(w, req)
	} else {
		r.controllers.Auth.RenderRegisterPage(w, req)
	}
}
