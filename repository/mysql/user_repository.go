package mysql

import (
	"database/sql"
	"errors"
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
	id, err := result.LastInsertId() //返回最后插入行的自增 ID
	if err != nil {
		return err
	}

	user.ID = int(id)           //将自增ID赋值给用户ID
	user.CreatedAt = time.Now() //将当前时间赋值给用户创建时间

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

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}

	query := `
		SELECT id, username, password, email, role, created_at 
		FROM users 
		WHERE username = ?
	`

	err := r.db.QueryRow(query, username).Scan(
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

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, password, email, role, created_at 
		FROM users 
		WHERE email = ?
	`
	err := r.db.QueryRow(query, email).Scan(
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

// GetAll 获取所有用户
func (r *userRepository) GetAll() ([]*models.User, error) {
	query := `
		SELECT id, username, email, role, created_at 
		FROM users 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// Update 更新用户信息
func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET username = ?, email = ?, role = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		user.Username,
		user.Email,
		user.Role,
		user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil

}

// Delete 删除用户
func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("用户不存在")
	}

	return nil
}

// Exists 检查用户是否存在
func (r *userRepository) Exists(username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = ?`

	err := r.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否已被使用
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Count 获取用户总数
func (r *userRepository) Count() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users`

	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
