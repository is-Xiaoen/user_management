package interfaces

import "user-management-system/models"

// UserRepository 定义用户数据访问接口
type UserRepository interface {
	//Creat 创建用户
	Create(user *models.User) error

	//GetByID 根据ID获取用户
	GetByID(id int) (*models.User, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(username string) (*models.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(email string) (*models.User, error)

	// GetAll 获取所有用户
	GetAll() ([]*models.User, error)

	// Update 更新用户信息
	Update(user *models.User) error

	// Delete 删除用户
	Delete(id int) error

	// Exists 检查用户是否存在
	Exists(username string) (bool, error)

	// ExistsByEmail 检查邮箱是否已被使用
	ExistsByEmail(email string) (bool, error)

	// Count 获取用户总数
	Count() (int64, error)
}
