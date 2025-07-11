package app

import (
	"database/sql"
	"time"

	"user-management-system/repository/interfaces"
	"user-management-system/repository/mysql"
	"user-management-system/services"
	"user-management-system/session"
)

// App 应用程序容器，管理所有依赖
type App struct {
	DB             *sql.DB
	UserRepository interfaces.UserRepository
	UserService    services.UserService
	SessionManager *session.Manager
}

// NewApp 创建应用实例
func NewApp(db *sql.DB, sessionCookieName string, sessionMaxLifetime time.Duration) *App {
	// 创建唯一的Repository实例
	userRepo := mysql.NewUserRepository(db)

	// 创建服务层
	service := services.NewService(&services.ServiceDependencies{
		DB:             db,
		UserRepository: userRepo,
	})

	// 创建会话管理器
	sessionManager := session.NewManager(sessionCookieName, sessionMaxLifetime)

	// 启动会话GC
	go sessionManager.GC()

	return &App{
		DB:             db,
		UserRepository: userRepo,
		UserService:    service.UserService,
		SessionManager: sessionManager,
	}
}
