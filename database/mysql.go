package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	
	"user-management-system/config"
	_ "github.com/go-sql-driver/mysql"
)

// DB 全局数据库连接实例
var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := config.GetConfig()
	
	// 构建DSN，添加超时参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&timeout=5s&readTimeout=5s&writeTimeout=5s",
   	 cfg.DBUser,
    	cfg.DBPassword,
    	cfg.DBHost,
    	cfg.DBPort,
    	cfg.DBName,
	)
	
	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %w", err)
	}
	
	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}
	
	DB = db
	log.Println("数据库连接成功")
	
	// 创建表（如果不存在）
	if err := createTables(db); err != nil {
		return fmt.Errorf("创建表失败: %w", err)
	}
	
	
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		} else {
			log.Println("数据库连接已关闭")
		}
	}
}

// createTables 创建必要的数据库表
func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		role VARCHAR(20) DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_username (username),
		INDEX idx_email (email),
		INDEX idx_role (role)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	
	_, err := db.Exec(query)
	return err
}



// GetDB 获取数据库连接实例
func GetDB() *sql.DB {
	return DB
}