package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management-system/config"
	"user-management-system/errors"

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

	// 获取配置
	cfg := config.GetConfig()

	// 创建服务器
	server := &http.Server{
		Addr:         "0.0.0.0:" + cfg.ServerPort,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      errors.RecoverMiddleware(http.DefaultServeMux), // 添加错误恢复中间件
	}

	// 创建通道监听终止信号
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		logger.Info("服务器启动在 http://localhost:%s", cfg.ServerPort)
		log.Printf("服务器启动在 http://localhost:%s\n", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("服务器启动失败: %v", err)
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待终止信号
	<-done
	logger.Info("收到关闭信号，服务器正在关闭...")
	log.Println("服务器正在关闭...")

	// 优雅关闭（给予5秒时间完成正在处理的请求）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败: %v", err)
		log.Printf("服务器关闭失败: %v", err)
	}

	logger.Info("服务器已停止")
	log.Println("服务器已停止")
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
			//fmt.Println("main中成功执行")
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
