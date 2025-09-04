package dao

import (
	"github.com/Orange-Health/citadel/apps/users/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetAllPathologists(filters structures.UserDetailDbFilters) ([]commonModels.User, *commonStructures.CommonError)
	GetUser(id uint) (commonModels.User, *commonStructures.CommonError)
	GetAllUsers() ([]commonModels.User, *commonStructures.CommonError)
	GetUsersByIds(userIds []uint) ([]commonModels.User, *commonStructures.CommonError)
	GetUserByEmail(email string) (commonModels.User, *commonStructures.CommonError)
	CreateUser(user commonModels.User) (commonModels.User, *commonStructures.CommonError)
	UpdateUser(user structures.UserDetail, updatedById uint) (commonModels.User, *commonStructures.CommonError)
	DeleteUser(userId, deletedById uint) *commonStructures.CommonError
}

func (userDao *UserDao) GetAllPathologists(filters structures.UserDetailDbFilters) ([]commonModels.User,
	*commonStructures.CommonError) {

	userTypes := []string{commonConstants.USER_TYPE_PATHOLOGIST, commonConstants.USER_TYPE_SUPER_ADMIN}

	users := []commonModels.User{}
	err := userDao.Db.Where("user_type IN (?)", userTypes).Limit(filters.Limit + 1).Offset(filters.Offset).Find(&users).Error
	if err != nil {
		return []commonModels.User{}, commonUtils.HandleORMError(err)
	}
	return users, nil
}

func (userDao *UserDao) GetUser(id uint) (commonModels.User, *commonStructures.CommonError) {

	user := commonModels.User{}
	err := userDao.Db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return commonModels.User{}, commonUtils.HandleORMError(err)
	}
	return user, nil
}

func (userDao *UserDao) GetUsersByIds(userIds []uint) ([]commonModels.User, *commonStructures.CommonError) {

	users := []commonModels.User{}
	err := userDao.Db.Where("id IN (?)", userIds).Find(&users).Error
	if err != nil {
		return []commonModels.User{}, commonUtils.HandleORMError(err)
	}

	return users, nil
}

func (userDao *UserDao) GetUserByEmail(email string) (commonModels.User, *commonStructures.CommonError) {

	user := commonModels.User{}
	err := userDao.Db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return commonModels.User{}, commonUtils.HandleORMError(err)
	}
	return user, nil
}

func (userDao *UserDao) GetAllUsers() ([]commonModels.User, *commonStructures.CommonError) {
	users := []commonModels.User{}
	err := userDao.Db.Find(&users).Error
	if err != nil {
		return []commonModels.User{}, commonUtils.HandleORMError(err)
	}
	return users, nil
}

func (userDao *UserDao) CreateUser(user commonModels.User) (commonModels.User, *commonStructures.CommonError) {

	err := userDao.Db.Create(&user).Error
	if err != nil {
		return commonModels.User{}, commonUtils.HandleORMError(err)
	}
	return user, nil
}

func (userDao *UserDao) UpdateUser(userDetail structures.UserDetail, updatedById uint) (commonModels.User,
	*commonStructures.CommonError) {

	updates := map[string]interface{}{
		"user_name":  userDetail.UserName,
		"email":      userDetail.Email,
		"user_type":  userDetail.UserType,
		"updated_by": updatedById,
	}
	err := userDao.Db.Model(&commonModels.User{}).Where("id = ?", userDetail.Id).Updates(updates).Error
	if err != nil {
		return commonModels.User{}, commonUtils.HandleORMError(err)
	}
	user := commonModels.User{}
	err = userDao.Db.Where("id = ?", userDetail.Id).First(&user).Error
	if err != nil {
		return commonModels.User{}, commonUtils.HandleORMError(err)
	}
	return user, nil
}

func (userDao *UserDao) DeleteUser(userId, deletedById uint) *commonStructures.CommonError {
	updates := map[string]interface{}{
		"updated_by": deletedById,
		"updated_at": commonUtils.GetCurrentTime(),
		"deleted_by": deletedById,
		"deleted_at": commonUtils.GetCurrentTime(),
	}

	err := userDao.Db.Model(&commonModels.User{}).Where("id = ?", userId).Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}
