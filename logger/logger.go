package logger // 定义包名为 logger

import (
	"fmt"        // 用于格式化字符串
	"log"        // 核心日志记录功能
	"os"         // 操作系统接口，用于文件操作
	"path/filepath" // 路径操作
	"sync"       // 并发原语，如互斥锁和 Once
	"time"       // 时间操作，用于日期格式化和计算时间间隔
)

var (
	fileLogger *Logger  // 全局的日志记录器实例
	once       sync.Once // 确保 Init 函数只被调用一次
)

// Logger 日志记录器结构体
type Logger struct {
	mu       sync.Mutex // 互斥锁，用于保护对文件和 logger 的并发访问
	file     *os.File   // 当前打开的日志文件句柄
	logger   *log.Logger // 标准库的 log.Logger 实例
	filePath string     // 当前日志文件的完整路径
}

// Init 初始化日志记录器
// logDir: 日志文件存储的目录
func Init(logDir string) error {
	var err error // 用于捕获 Init 过程中可能发生的错误
	once.Do(func() { // 确保 Init 函数内的代码只执行一次
		// 确保日志目录存在，如果不存在则创建
		// 0755 是目录的权限模式：所有者读写执行，组用户和其他用户只读执行
		if err = os.MkdirAll(logDir, 0755); err != nil {
			return // 如果创建目录失败，直接返回错误
		}

		// 生成日志文件名，格式为 "app_YYYY-MM-DD.log"
		fileName := fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02"))
		// 拼接日志目录和文件名，得到完整的日志文件路径
		filePath := filepath.Join(logDir, fileName)

		// 打开或创建日志文件
		// os.O_CREATE: 如果文件不存在则创建
		// os.O_APPEND: 写入时追加到文件末尾
		// os.O_WRONLY: 只写模式
		// 0644 是文件的权限模式：所有者读写，组用户和其他用户只读
		file, openErr := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if openErr != nil {
			err = openErr // 如果打开文件失败，记录错误
			return
		}

		// 初始化全局的 fileLogger 实例
		fileLogger = &Logger{
			file:     file,                                 // 设置文件句柄
			logger:   log.New(file, "", log.LstdFlags), // 创建一个新的 log.Logger，输出到文件，不带前缀，包含标准日期时间
			filePath: filePath,                             // 设置文件路径
		}

		// 启动一个 goroutine，用于每天轮转日志文件
		go fileLogger.rotateDaily(logDir)
	})

	return err // 返回初始化过程中可能发生的错误
}

// rotateDaily 每天轮转日志文件
// logDir: 日志文件存储的目录
func (l *Logger) rotateDaily(logDir string) {
	for { // 无限循环，持续进行日志轮转
		// 计算到明天凌晨的时间点
		now := time.Now()
		// 获取明天凌晨 0 点 0 分 0 秒的时间
		tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		// 计算从现在到明天凌晨的时间间隔
		duration := tomorrow.Sub(now)

		// 等待到明天凌晨，这个 goroutine 会在这里阻塞
		time.Sleep(duration)

		// 达到明天凌晨后，开始轮转操作
		l.mu.Lock() // 获取互斥锁，防止在轮转过程中有其他协程写入日志

		// 关闭旧的日志文件
		if l.file != nil {
			l.file.Close()
		}

		// 创建新的日志文件
		fileName := fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02"))
		filePath := filepath.Join(logDir, fileName)
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			// 如果创建新文件失败，打印错误信息（这里用标准的 log 包打印到 stderr）
			log.Printf("创建新日志文件失败: %v", err)
			l.mu.Unlock() // 释放锁
			continue      // 继续下一次循环，等待下一个轮转周期
		}

		// 更新 Logger 结构体中的文件句柄、logger 实例和文件路径
		l.file = file
		l.logger = log.New(file, "", log.LstdFlags)
		l.filePath = filePath

		l.mu.Unlock() // 释放互斥锁
	}
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	if fileLogger != nil { // 检查日志记录器是否已初始化
		fileLogger.mu.Lock()         // 获取互斥锁，确保并发安全
		defer fileLogger.mu.Unlock() // 函数执行完毕后释放锁
		fileLogger.logger.Printf("[INFO] "+format, v...) // 写入信息日志
	}
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	if fileLogger != nil {
		fileLogger.mu.Lock()
		defer fileLogger.mu.Unlock()
		fileLogger.logger.Printf("[ERROR] "+format, v...) // 写入错误日志
	}
}

// Warning 记录警告日志
func Warning(format string, v ...interface{}) {
	if fileLogger != nil {
		fileLogger.mu.Lock()
		defer fileLogger.mu.Unlock()
		fileLogger.logger.Printf("[WARNING] "+format, v...) // 写入警告日志
	}
}

// UserAction 记录用户操作日志
// username: 用户名
// action: 操作内容
// details: 操作详情
// success: 操作是否成功
func UserAction(username, action, details string, success bool) {
	status := "成功"
	if !success {
		status = "失败"
	}
	// 调用 Info 函数记录用户操作日志，格式化输出
	Info("用户操作 - 用户: %s, 操作: %s, 详情: %s, 结果: %s", username, action, details, status)
}

// UserActionWithError 记录用户操作日志（包含错误详情）
func UserActionWithError(username, action, details string, err error) {
    if err != nil {
        // 如果有错误，记录详细错误信息
        Error("用户操作失败 - 用户: %s, 操作: %s, 详情: %s, 错误: %v", 
            username, action, details, err)
    } else {
        // 成功的情况
        Info("用户操作成功 - 用户: %s, 操作: %s, 详情: %s", 
            username, action, details)
    }
}

// Close 关闭日志文件
func Close() {
	// 在程序退出时调用，确保关闭打开的日志文件
	if fileLogger != nil && fileLogger.file != nil {
		fileLogger.file.Close()
	}
}