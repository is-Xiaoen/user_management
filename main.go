package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user-management-system/app"
	"user-management-system/config"
	"user-management-system/database"
	"user-management-system/errors"
	"user-management-system/logger"
	"user-management-system/router"
)

func main() {
	// 初始化日志记录器
	if err := logger.Init("logs"); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	defer logger.Close()

	logger.Info("应用程序启动中...")

	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		logger.Error("数据库初始化失败: %v", err)
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.CloseDB()

	// 创建应用实例（统一管理所有依赖）
	application := app.NewApp(
		database.GetDB(),
		"session_id",
		2*time.Hour,
	)

	// 创建路由器
	r := router.NewRouter(application)
	handler := r.Setup()

	// 获取配置
	cfg := config.GetConfig()

	// 创建服务器
	server := &http.Server{
		Addr:         "0.0.0.0:" + cfg.ServerPort,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      errors.RecoverMiddleware(handler),
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
