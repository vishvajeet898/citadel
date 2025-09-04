package controller

import "github.com/Orange-Health/citadel/apps/users/service"

type User struct {
	UserService service.UserServiceInterface
}

func InitUserController() *User {
	return &User{
		UserService: service.InitializeUserService(),
	}
}
