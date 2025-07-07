package interfaces

import "user-management-system/models"

// UserRepository 定义用户数据访问接口
type UserRepository interface{
	//Creat 创建用户
	Create(user *models.User) error 
	
	//GetByID 根据ID获取用户
	GetByID(id int) (*models.User, error)
}

