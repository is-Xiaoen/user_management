package services

import "user-management-system/repository/interfaces"

//UserService 用户服务接口
type UserService interface{
	
}

// userServiceImpl 是 UserService 接口的具体实现
type userServiceImpl struct{
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) UserService{
	return &userServiceImpl{
		userRepo: userRepo,
	}
}