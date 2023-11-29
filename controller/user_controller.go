package controller

import (
	"TeamSyncMessenger-Backend/helper"
	"TeamSyncMessenger-Backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUsers(c *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *userController {
	return &userController{
		userService: userService,
	}
}

func (us *userController) GetUsers(c *gin.Context) {
	users, err := us.userService.GetUsers()
	if err != nil {
		res := helper.BuildErrorResponse("사용자 목록을 찾을 수 없습니다.", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	res := helper.BuildResponse(true, "사용자 목록", users)
	c.JSON(200, res)
}
