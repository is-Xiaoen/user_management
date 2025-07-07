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
