package services

import (
	"fmt"
	"user-management-system/errors"
	"user-management-system/models"
	"user-management-system/repository/interfaces"
)

// UserService 用户服务接口
type UserService interface {
	// 用户认证相关
	RegisterUser(username, password, email string) error
	AuthenticateUser(username, password string) (*models.User, error)

	//用户管理相关
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUser(id int, email, role string) error
	DeleteUser(id int) error

	//权限检查
	IsAdmin(user *models.User) bool

	//统计相关
	GetUserStats() (map[string]interface{}, error)
}

// userServiceImpl 是 UserService 接口的具体实现
type userServiceImpl struct {
	userRepo interfaces.UserRepository
}

// NewUserService 创建一个新的用户服务实例
func NewUserService(userRepo interfaces.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

// RegisterUser 注册一个新用户
func (s *userServiceImpl) RegisterUser(username, password, email string) error {
	//验证输入
	if username == "" {
		return errors.NewValidationError("username", "用户名不能为空")
	}
	if len(username) < 3 || len(username) > 20 {
		return errors.NewValidationError("username", "用户名长度必须在3到20个字符之间")
	}
	if password == "" {
		return errors.NewValidationError("password", "密码不能为空")
	}
	if len(password) < 6 || len(password) > 20 {
		return errors.NewValidationError("password", "密码长度必须在6到20个字符之间")
	}
	if email == "" {
		return errors.NewValidationError("email", "邮箱不能为空")
	}
	// 检查用户名是否已存在
	exists, err := s.userRepo.Exists(username)
	if err != nil {
		return errors.NewInternalError(err)
	}
	if exists {
		return errors.NewConflictError("用户名已存在")
	}

	//检查邮箱是否已经被使用
	emailExists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("检查邮箱失败: %w", err))
	}

	if emailExists {
		return errors.NewConflictError("邮箱已被使用")
	}

	//创建新用户
	user := &models.User{
		Username: username,
		Email:    email,
		Role:     "user", //默认角色
	}

	//设置密码(使用bcrypt加密)
	if err := user.SetPassword(password); err != nil {
		return errors.NewInternalError(fmt.Errorf("设置密码失败: %w", err))
	}

	//保存到数据库
	if err := s.userRepo.Create(user); err != nil {
		return errors.NewInternalError(fmt.Errorf("保存用户失败: %w", err))
	}
	return nil
}

// AuthenticateUser 用户认证
func (s *userServiceImpl) AuthenticateUser(username, password string) (*models.User, error) {
	//fmt.Println("进入到services__user_service.go 100")

	// 验证输入
	if username == "" || password == "" {
		return nil, errors.NewValidationError("", "用户名和密码不能为空")
	}

	//获取用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("获取用户失败: %w", err))
	}
	if user == nil {
		return nil, errors.NewUnauthorizedError("用户不存在")
	}
	if !user.CheckPassword(password) {
		return nil, errors.NewUnauthorizedError("密码错误")
	}
	return user, nil
}

// GetUserByID 通过ID获取用户
func (s *userServiceImpl) GetUserByID(id int) (*models.User, error) {
	if id <= 0 {
		return nil, errors.NewValidationError("id", "无效的用户ID")
	}
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("获取用户失败: %w", err))
	}
	if user == nil {
		return nil, errors.NewNotFoundError("用户不存在")
	}
	return user, nil
}

// GetUserByUsername 通过用户名获取用户
func (s *userServiceImpl) GetUserByUsername(username string) (*models.User, error) {
	if username == "" {
		return nil, errors.NewValidationError("username", "用户名不能为空")
	}

	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("获取用户失败: %w", err))
	}
	if user == nil {
		return nil, errors.NewNotFoundError("用户不存在")
	}
	return user, nil
}

// GetAllUsers 获取所有用户
func (s *userServiceImpl) GetAllUsers() ([]*models.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("获取用户列表失败: %w", err))
	}
	return users, nil
}

// UpdateUser 更新用户信息
func (s *userServiceImpl) UpdateUser(id int, email, role string) error {
	if id <= 0 {
		return errors.NewValidationError("id", "无效的用户ID")
	}
	if email == "" {
		return errors.NewValidationError("email", "邮箱不能为空")
	}
	if role != "user" && role != "admin" {
		return errors.NewValidationError("role", "无效的角色")
	}

	// 检查用户是否存在
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("查询用户失败: %w", err))
	}
	if existingUser == nil {
		return errors.NewNotFoundError("用户")
	}

	// 如果邮箱改变了，检查新邮箱是否已被使用
	if existingUser.Email != email {
		emailUser, err := s.userRepo.GetByEmail(email)
		if err != nil {
			return errors.NewInternalError(fmt.Errorf("检查邮箱失败: %w", err))
		}

		if emailUser != nil && emailUser.ID != id {
			return errors.NewConflictError("邮箱已被其他用户使用")
		}
	}

	//更新用户信息
	if err := s.userRepo.UpdateEmailAndRole(id, email, role); err != nil {
		return errors.NewInternalError(fmt.Errorf("更新用户信息失败: %w", err))
	}

	return nil
}

// DeleteUser 删除用户
func (s *userServiceImpl) DeleteUser(id int) error {
	//验证输入
	if id <= 0 {
		return errors.NewValidationError("id", "无效的用户ID")
	}

	//检查用户存在
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("查询用户失败: %w", err))
	}
	if user == nil {
		return errors.NewNotFoundError("用户")
	}

	//防止删除最后一个管理员
	if user.Role == "admin" {
		adminCount, err := s.userRepo.CountByRole("admin")
		if err != nil {
			return errors.NewInternalError(fmt.Errorf("查询管理员数量失败: %w", err))
		}
		if adminCount <= 1 {
			return errors.NewForbiddenError("不能删除最后一个管理员")
		}
	}

	//删除用户
	if err := s.userRepo.Delete(id); err != nil {
		return errors.NewInternalError(fmt.Errorf("删除用户失败: %w", err))
	}

	return nil
}

// IsAdmin 检查用户是否为管理员
func (s *userServiceImpl) IsAdmin(user *models.User) bool {
	return user != nil && user.IsAdmin()
}

// GetUserStats 获取用户统计信息
func (s *userServiceImpl) GetUserStats() (map[string]interface{}, error) {
	totalCount, err := s.userRepo.Count()
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("获取用户总数失败: %w", err))
	}
	adminCount, err := s.userRepo.CountByRole("admin")
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("获取管理员数量失败: %w", err))
	}
	userCount, err := s.userRepo.CountByRole("user")
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("获取普通用户数量失败: %w", err))
	}

	stats := map[string]interface{}{
		"total": totalCount,
		"admin": adminCount,
		"user":  userCount,
	}

	return stats, nil
}
