package models

import "time"

//User 表示用户模型, 映射数据库中的users表
type User struct {
	ID        int       `json:"id"`         // 用户 ID
	Username  string    `json:"username"`   // 用户名
	Password  string    `json:"-"`          // 密码（JSON序列化时忽略）
	Email     string    `json:"email"`      // 邮箱
	Role      string    `json:"role"`       // 角色（user/admin）
	CreatedAt time.Time `json:"created_at"` // 创建时间
}


