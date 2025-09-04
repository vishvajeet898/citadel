package service

import (
	"net/http"

	mapper "github.com/Orange-Health/citadel/apps/users/mapper"
	"github.com/Orange-Health/citadel/apps/users/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type UserServiceInterface interface {
	GetAllPathologists(filters structures.UserDetailDbFilters) (
		[]structures.UserDetail, *commonStructures.CommonError)
	GetAllPathologistsModels() ([]commonModels.User, *commonStructures.CommonError)
	GetUser(userID uint) (structures.UserDetail, *commonStructures.CommonError)
	GetUserModel(userID uint) (commonModels.User, *commonStructures.CommonError)
	GetUsersByIds(userIds []uint) ([]commonModels.User, *commonStructures.CommonError)
	GetUserByEmail(email string) (structures.UserDetail, *commonStructures.CommonError)
	GetUserIdNameMap() (map[uint]string, *commonStructures.CommonError)

	CreateUser(userDetail structures.UserDetail, createdById uint) (structures.UserDetail, *commonStructures.CommonError)
	UpdateUser(userDetail structures.UserDetail, updatedById uint) (structures.UserDetail, *commonStructures.CommonError)
	DeleteUser(userID uint, deletedById uint) *commonStructures.CommonError
}

func validateUserType(userType string) *commonStructures.CommonError {
	if !commonUtils.SliceContainsString(commonConstants.USER_TYPES, userType) {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.INVALID_USER_TYPE,
		}
	}

	return nil
}

func validateUserCreationUpdation(userDetail structures.UserDetail) *commonStructures.CommonError {
	if userDetail.UserName == "" {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.USER_NAME_NOT_FOUND,
		}
	}

	if userDetail.Email == "" {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.EMAIL_NOT_FOUND,
		}
	}

	cErr := validateUserType(userDetail.UserType)
	if cErr != nil {
		return cErr
	}

	return nil
}

func (userService *UserService) GetUser(userID uint) (structures.UserDetail, *commonStructures.CommonError) {
	user, err := userService.UserDao.GetUser(userID)
	if err != nil {
		return structures.UserDetail{}, err
	}

	return mapper.MapUserDetail(user), nil
}

func (userService *UserService) GetUserModel(userID uint) (commonModels.User, *commonStructures.CommonError) {
	user, err := userService.UserDao.GetUser(userID)
	if err != nil {
		return commonModels.User{}, err
	}

	return user, nil
}

func (userService *UserService) GetUsersByIds(userIds []uint) ([]commonModels.User, *commonStructures.CommonError) {
	users, err := userService.UserDao.GetUsersByIds(userIds)
	if err != nil {
		return []commonModels.User{}, err
	}

	return users, nil
}

func (userService *UserService) GetAllPathologists(filters structures.UserDetailDbFilters) (
	[]structures.UserDetail, *commonStructures.CommonError) {
	users, err := userService.UserDao.GetAllPathologists(filters)
	if err != nil {
		return []structures.UserDetail{}, err
	}

	return mapper.MapUserDetails(users), nil
}

func (userService *UserService) GetAllPathologistsModels() ([]commonModels.User, *commonStructures.CommonError) {
	users, err := userService.UserDao.GetAllPathologists(structures.UserDetailDbFilters{
		Limit: 100,
	})
	if err != nil {
		return []commonModels.User{}, err
	}

	return users, nil
}

func (userService *UserService) GetUserByEmail(email string) (structures.UserDetail, *commonStructures.CommonError) {
	user, cErr := userService.UserDao.GetUserByEmail(email)
	if cErr != nil {
		return structures.UserDetail{}, cErr
	}

	return mapper.MapUserDetail(user), nil
}

func (userService *UserService) GetUserIdNameMap() (map[uint]string, *commonStructures.CommonError) {
	userIdNameMap := make(map[uint]string)
	users, cErr := userService.UserDao.GetAllUsers()
	if cErr != nil {
		return userIdNameMap, cErr
	}

	for _, user := range users {
		userIdNameMap[user.Id] = user.UserName
	}

	return userIdNameMap, nil
}

func (userService *UserService) CreateUser(userDetail structures.UserDetail, createdById uint) (
	structures.UserDetail, *commonStructures.CommonError) {

	user := mapper.MapUser(userDetail)
	cErr := validateUserCreationUpdation(userDetail)
	if cErr != nil {
		return structures.UserDetail{}, cErr
	}

	dbUser, cErr := userService.UserDao.GetUserByEmail(userDetail.Email)
	if cErr != nil {
		return structures.UserDetail{}, cErr
	}
	if dbUser.Id != 0 {
		return structures.UserDetail{}, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.USER_ALREADY_EXISTS,
		}
	}

	user.CreatedBy = createdById
	user.UpdatedBy = createdById
	user, cErr = userService.UserDao.CreateUser(user)
	if cErr != nil {
		return structures.UserDetail{}, cErr
	}

	return mapper.MapUserDetail(user), nil
}

func (userService *UserService) UpdateUser(userDetail structures.UserDetail, updatedById uint) (
	structures.UserDetail, *commonStructures.CommonError) {

	cErr := validateUserCreationUpdation(userDetail)
	if cErr != nil {
		return structures.UserDetail{}, cErr
	}

	user, cErr := userService.UserDao.UpdateUser(userDetail, updatedById)
	if cErr != nil {
		return structures.UserDetail{}, cErr
	}

	return mapper.MapUserDetail(user), nil
}

func (userService *UserService) DeleteUser(userId, deletedById uint) *commonStructures.CommonError {
	cErr := userService.UserDao.DeleteUser(userId, deletedById)
	if cErr != nil {
		return cErr
	}

	return nil
}
