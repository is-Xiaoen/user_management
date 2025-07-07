package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User 表示用户模型, 映射数据库中的users表
type User struct {
	ID        int       `json:"id"`         // 用户 ID
	Username  string    `json:"username"`   // 用户名
	Password  string    `json:"-"`          // 密码（JSON序列化时忽略）
	Email     string    `json:"email"`      // 邮箱
	Role      string    `json:"role"`       // 角色（user/admin）
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// CheckPassword
func (u *User) CheckPassword(password string) bool {
	// 明文比较
	if u.Password == password {
		return true
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// SetPassword 设置用户密码（自动进行哈希）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Validate 验证用户数据
func (u *User) Validate() error {
	// 这里可以添加更多的验证逻辑
	return nil
}
