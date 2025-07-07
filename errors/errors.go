package errors

import (
	"fmt"
	"log"
	"net/http"
)

// ErrorType 错误类型
type ErrorType int

const (
	//ValidationError 验证错误
	ValidationError ErrorType = iota
	// NotFoundError 资源不存在
	NotFoundError
	//UnauthorizedError 未授权
	UnauthorizedError
	//ForbiddenError 禁止访问
	ForbiddenError
	// ConflictError 冲突错误 (例如用户名已存在)
	ConflictError
	// InternalError 内部服务器错误
	InternalError
)

// AppError 应用错误结构
type AppError struct {
	Type     ErrorType
	Message  string      //用户友好的错误信息(展示给用户的)
	Internal error       //内部错误, 用于记录日志(开发人员使用)
	Field    string      //相关字段(用于表单验证) 说明哪个字段出错
	Data     interface{} //额外数据(用于传递数据)
}

// Error 实现error接口
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError 创建新的应用错误
func NewAppError(errType ErrorType, message string, internal error) *AppError {
	return &AppError{
		Type:     errType,
		Message:  message,
		Internal: internal,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(field, message string) *AppError {
	return &AppError{
		Type:    ValidationError,
		Message: message,
		Field:   field,
	}
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Type:    NotFoundError,
		Message: fmt.Sprintf("%s不存在", resource),
	}
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "请先登录"
	}
	return &AppError{
		Type:    UnauthorizedError,
		Message: message,
	}
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "权限不足"
	}
	return &AppError{
		Type:    ForbiddenError,
		Message: message,
	}
}

// NewConflictError 创建冲突错误
func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ConflictError,
		Message: message,
	}
}

// NewInternalError 创建内部错误
func NewInternalError(internal error) *AppError {
	return &AppError{
		Type:     InternalError,
		Message:  "服务器内部错误，请稍后重试"+internal.Error(),
		Internal: internal,
	}
}

// HTTPStatusCode 获取HTTP状态码
func (e *AppError) HTTPStatusCode() int {
	switch e.Type {
	case ValidationError:
		return http.StatusBadRequest // 400
	case NotFoundError:
		return http.StatusNotFound // 404
	case UnauthorizedError:
		return http.StatusUnauthorized
	case ForbiddenError:
		return http.StatusForbidden
	case ConflictError:
		return http.StatusConflict
	case InternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError // 500
	}
}

// LogError 记录错误日志
func (e *AppError) LogError() {
	if e.Internal != nil {
		log.Printf("[错误] 类型: %s, 消息: %s, 内部错误: %v",
			e.TypeString(), e.Message, e.Internal)
	} else {
		log.Printf("[错误] 类型: %s, 消息: %s",
			e.TypeString(), e.Message)
	}
}

// TypeString 获取错误类型字符串
func (e *AppError) TypeString() string {
	switch e.Type {
	case ValidationError:
		return "验证错误"
	case NotFoundError:
		return "资源不存在"
	case UnauthorizedError:
		return "未授权"
	case ForbiddenError:
		return "禁止访问"
	case ConflictError:
		return "冲突"
	case InternalError:
		return "内部错误"
	default:
		return "未知错误"
	}
}

// IsAppError 检查是否为应用错误
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// HandleError 统一处理错误响应
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	// 检查是否为应用错误
	appErr, ok := IsAppError(err)
	if !ok {
		// 如果不是应用错误，创建一个内部错误
		appErr = NewInternalError(err)
	}

	// 记录错误日志
	appErr.LogError()

	// 根据请求类型返回不同格式的响应
	if isAPIRequest(r) {
		// API请求返回JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.HTTPStatusCode())
		fmt.Fprintf(w, `{"error":"%s"}`, appErr.Message)
	} else {
		// 普通请求返回HTML错误页面
		http.Error(w, appErr.Message, appErr.HTTPStatusCode())
	}
}

// isAPIRequest 判断是否为API请求
func isAPIRequest(r *http.Request) bool {
	return len(r.URL.Path) > 4 && r.URL.Path[:4] == "/api"
}

// RecoverMiddleware 恢复中间件，捕获panic
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v", err)
				appErr := NewInternalError(fmt.Errorf("panic: %v", err))
				HandleError(w, r, appErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
