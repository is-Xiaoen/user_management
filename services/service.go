package services

import (
	"database/sql"
	"user-management-system/repository/interfaces"
	"user-management-system/repository/mysql"
)

// Service 是所有服务的集合，用于统一管理服务实例
type Service struct {
	UserService UserService
}

// ServiceDependencies 服务依赖项
type ServiceDependencies struct {
	DB             *sql.DB
	UserRepository interfaces.UserRepository
}

// NewService  创建一个新的服务集合实例
func NewService(deps *ServiceDependencies) *Service {
	// 如果没有提供UserRepository，使用默认的MySQL实现
	if deps.UserRepository == nil {
		deps.UserRepository = mysql.NewUserRepository(deps.DB)
	}
	return &Service{
		UserService: NewUserService(deps.UserRepository),
	}
}

// NewServiceWithDB 使用数据库连接创建服务（便捷方法）
func NewServiceWithDB(db *sql.DB) *Service {
	return NewService(&ServiceDependencies{
		DB:             db,
		UserRepository: mysql.NewUserRepository(db),
	})
}
