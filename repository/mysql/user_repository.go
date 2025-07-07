package mysql

import (
	"database/sql"
	"time"
	"user-management-system/models"
	"user-management-system/repository/interfaces"
)

// userRepository MySQL实现的用户仓库
type userRepository struct {
	db *sql.DB
}

// NewUserRepository 创建MySQL用户仓库实例
func NewUserRepository(db *sql.DB) interfaces.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create 创建新用户
func (r *userRepository) Create(user *models.User) error {
	//防止 SQL 注入攻击
	query := `
		INSERT INTO users (username, password, email, role, created_at) 
		VALUES (?, ?, ?, ?, ?)
	`
	
	result, err := r.db.Exec(query, 
		user.Username, 
		user.Password, 
		user.Email, 
		user.Role, 
		time.Now(),
	)
	
	if err != nil {
		return err
	}
	
	// 获取插入的ID
	id, err := result.LastInsertId()//返回最后插入行的自增 ID    
	if err != nil {
		return err
	}
	
	user.ID = int(id)//将自增ID赋值给用户ID
	user.CreatedAt = time.Now()//将当前时间赋值给用户创建时间
	
	return nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	
	query := `
		SELECT id, username, password, email, role, created_at 
		FROM users 
		WHERE id = ?
	`
	
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return user, nil
}