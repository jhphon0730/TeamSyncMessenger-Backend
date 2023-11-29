package controller

import (
	"TeamSyncMessenger-Backend/service"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	TestUser(c *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *userController {
	return &userController{
		userService: userService,
	}
}

func (us *userController) TestUser(c *gin.Context) {
	userName := us.userService.TestUser()
	c.JSON(200, userName)
}
