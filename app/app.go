package app

import (
	"database/sql"
	"time"

	"user-management-system/session"
)

// App 应用程序容器，只管理全局共享的依赖
type App struct {
	DB             *sql.DB
	SessionManager *session.Manager
	// 移除了 UserRepository 和 UserService
}

// NewApp 创建应用实例
func NewApp(db *sql.DB, sessionCookieName string, sessionMaxLifetime time.Duration) *App {
	// 创建会话管理器
	sessionManager := session.NewManager(sessionCookieName, sessionMaxLifetime)

	// 启动会话GC
	go sessionManager.GC()

	return &App{
		DB:             db,
		SessionManager: sessionManager,
	}
}

// GetDB 获取数据库连接（供控制器使用）
func (a *App) GetDB() *sql.DB {
	return a.DB
}

// GetSessionManager 获取会话管理器（供控制器使用）
func (a *App) GetSessionManager() *session.Manager {
	return a.SessionManager
}
