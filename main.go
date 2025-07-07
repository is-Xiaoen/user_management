package main

import (
	"log"
	"net/http"
	"time"

	"user-management-system/controllers"
	"user-management-system/database"
	"user-management-system/logger"
	"user-management-system/middleware"
	"user-management-system/repository/mysql"
	"user-management-system/services"
	"user-management-system/session"
)

func main() {

	// 初始化日志记录器
	if err := logger.Init("logs"); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	defer logger.Close()

	logger.Info("应用程序启动中...")

	//初始化数据库连接
	if err := database.InitDB(); err != nil {
		log.Fatalf("无法初始化数据库: %v", err)
	}
	defer database.CloseDB()

	//初始化会话管理器 (会话有效期2小时)
	session.Init("session_id", 2*time.Hour)

	// 初始化Repository
	userRepo := mysql.NewUserRepository(database.GetDB())

	//session会话层和service服务层都有一个UserRepository实例,这两个实例不一样
	//repository模式是为了分离model层的数据访问层
	//而session和service层 中都需要查询数据库,所以这两个层中都有一个UserRepository实例

	// 初始化会话模块的Repository ,主要是为了根据ID查询完整的用户信息 session结构体中含用户ID字段
	session.InitSession(userRepo)

	//初始化服务层 (使用新的Repository模式)处理完整的用户业务逻辑，需要全套的CRUD操作
	service := services.NewServiceWithDB(database.GetDB())

	//初始化控制器和中间件
	controllers.InitControllers(service)
	middleware.InitMiddleware(service)

	//设置路由
	setupRoutes()
}

// 设置路由
func setupRoutes() {
	//创建一个带CSRF保护的处理器链
	cerfChain := session.CSRFMiddleware

	//静态文件服务
	//这段代码的作用是：让你的网站能够提供静态文件服务 意思是在浏览器
	//直接通过http://localhost:8080/static/css/style.css方式访问静态文件时可以正常显示出静态文件
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//主页路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		controllers.RenderHomePage(w, r)
	})

	//登录路由
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.HandleLogin(w, r)
		} else {
			controllers.RenderLoginPage(w, r)
		}
	})

	//注册路由
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.HandleRegister(w, r)
		} else {
			controllers.RenderRegisterPage(w, r)
		}
	})

	//登出路由
	http.HandleFunc("/logout", controllers.HandleLogout)

	//用户列表路由 (需要认证)
	http.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(controllers.RenderUsersPage)))

	//用户删除路由 (管理员权限 + CSRF保护)
	http.Handle("/users/delete", middleware.AuthMiddleware(
		cerfChain(
			http.HandlerFunc(controllers.HandleDeleteUser))))

	//用户更新路由 (管理员权限 + CSRF保护)
	http.Handle("/users/update", middleware.AuthMiddleware(
		cerfChain(
			http.HandlerFunc(controllers.HandleUpdateUser))))

	// API路由（用于返回JSON数据）
	http.HandleFunc("/api/users", controllers.HandleAPIUsers)
	http.HandleFunc("/api/users/stats", controllers.HandleAPIUserStats)

}
