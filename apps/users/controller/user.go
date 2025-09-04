package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/users/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		Get Users
// @Description	Get all users
// @Tags			users
// @Produce		json
// @Param			limit	query		int						false	"Limit"
// @Param			offset	query		int						false	"Offset"
// @Success		200		{object}	[]structures.UserDetail	"User Details"
// @Router			/api/v1/users [get]
func (userController *User) GetAllUsers(c *gin.Context) {

	filters := structures.UserDetailDbFilters{
		Limit:  commonUtils.ConvertStringToInt(c.Query("limit")),
		Offset: commonUtils.ConvertStringToInt(c.Query("offset")),
	}

	users, err := userController.UserService.GetAllPathologists(filters)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary		Get User
// @Description	Get a user by id
// @Tags			users
// @Produce		json
// @Param			id	path		int						true	"User ID"
// @Success		200	{object}	structures.UserDetail	"User Detail"
// @Router			/api/v1/users/{id} [get]
func (userController *User) GetUser(c *gin.Context) {
	// Get the user
	id := commonUtils.ConvertStringToUint(c.Param("id"))

	if id == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.INVALID_USER_ID_RECEIVED)
		return

	}

	user, err := userController.UserService.GetUser(id)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return

	}
	c.JSON(http.StatusOK, user)
}

// @Summary		Create User
// @Description	Create a user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			user	body		structures.UserDetail	true	"User Details"
// @Success		200		{object}	structures.UserDetail	"User Details"
// @Router			/api/v1/user [post]
func (userController *User) CreateUser(c *gin.Context) {
	createdById, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	if createdById == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.INVALID_USER_ID_RECEIVED)
		return
	}

	var userDetail structures.UserDetail

	if err := c.ShouldBindJSON(&userDetail); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := userController.UserService.CreateUser(userDetail, createdById)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return

	}
	c.JSON(http.StatusOK, user)
}

// @Summary		Update User
// @Description	Update a user by id
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id		path		int						true	"User ID"
// @Param			user	body		structures.UserDetail	true	"User Details"
// @Success		200		{object}	structures.UserDetail	"User Details"
// @Router			/api/v1/user/{id} [post]
func (userController *User) UpdateUser(c *gin.Context) {
	updatedById, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	if updatedById == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.INVALID_USER_ID_RECEIVED)
		return
	}
	var userDetail structures.UserDetail

	userID := commonUtils.ConvertStringToUint(c.Param("id"))
	if userID == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.INVALID_USER_ID_RECEIVED)
		return
	}

	if err := c.ShouldBindJSON(&userDetail); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userDetail.Id = userID

	user, err := userController.UserService.UpdateUser(userDetail, updatedById)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (userController *User) DeleteUser(c *gin.Context) {
	updatedById, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	if updatedById == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.INVALID_USER_ID_RECEIVED)
		return
	}

	userId := commonUtils.ConvertStringToUint(c.Param("id"))
	if userId == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.INVALID_USER_ID_RECEIVED)
		return
	}

	err := userController.UserService.DeleteUser(userId, updatedById)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

// @Summary		Get User By Email
// @Description	Get a user by email
// @Tags			users
// @Produce		json
// @Param			email	query		string					true	"Email"
// @Success		200		{object}	structures.UserDetail	"User Details"
// @Router			/api/v1/users/get-user-by-email [get]
func (userController *User) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")

	if email == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.INVALID_EMAIL)
		return
	}

	user, err := userController.UserService.GetUserByEmail(email)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, user)
}
