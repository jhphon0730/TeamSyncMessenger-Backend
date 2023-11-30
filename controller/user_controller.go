package controller

import (
	"TeamSyncMessenger-Backend/DTO"
	"TeamSyncMessenger-Backend/helper"
	"TeamSyncMessenger-Backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUsers(c *gin.Context)
	RegisterUser(c *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *userController {
	return &userController{
		userService: userService,
	}
}

func (uc *userController) GetUsers(c *gin.Context) {
	users, err := uc.userService.GetUsers()
	if err != nil {
		res := helper.BuildErrorResponse("사용자 목록을 찾을 수 없습니다.", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	res := helper.BuildResponse(true, "사용자 목록", users)
	c.JSON(200, res)
}

func (uc *userController) RegisterUser(c *gin.Context) {
	var registerUserDTO DTO.RegisterUserDTO
	if err := c.ShouldBindJSON(&registerUserDTO); err != nil {
		res := helper.BuildErrorResponse(`입력된 데이터가 정상이 아닙니다.`, err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if registerUserDTO.Username == "" || registerUserDTO.Password == "" || registerUserDTO.Email == "" {
		res := helper.BuildErrorResponse(`입력된 데이터가 없습니다.`, "", helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// GetUserByUsername 함수는 매개변수로 넘긴 Username으로 등록 된 사용자가 없으면 err 를 반환함
	validUser, _ := uc.userService.GetUserByUsername(registerUserDTO.Username)
	if validUser.ID != 0 {
		res := helper.BuildErrorResponse(`동일한 사용자 정보가 존재합니다.`, "", helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	user, err := uc.userService.CreateUser(registerUserDTO)
	if err != nil {
		res := helper.BuildErrorResponse(`정상적인 데이터가 아닙니다.`, err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, `사용자 회원가입에 성공하였습니다.`, user)
	c.JSON(201, res)
}
