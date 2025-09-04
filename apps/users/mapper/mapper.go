package mapper

import (
	"github.com/Orange-Health/citadel/apps/users/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapUser(user structures.UserDetail) commonModels.User {

	var userModel commonModels.User
	userModel.UserName = user.UserName
	userModel.UserType = user.UserType
	userModel.Email = user.Email

	return userModel
}

func MapUserDetail(user commonModels.User) structures.UserDetail {
	var userStruct structures.UserDetail

	userStruct.Id = user.Id
	userStruct.UserName = user.UserName
	userStruct.UserType = user.UserType
	userStruct.Email = user.Email

	return userStruct
}

func MapUserDetails(users []commonModels.User) []structures.UserDetail {
	uds := []structures.UserDetail{}
	for _, user := range users {
		uds = append(uds, MapUserDetail(user))
	}
	return uds
}
