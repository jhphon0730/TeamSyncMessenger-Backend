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
	LoginUser(c *gin.Context)
}

type userController struct {
	userService service.UserService
	authService service.AuthService
}

func NewUserController(userService service.UserService, authService service.AuthService) *userController {
	return &userController{
		userService: userService,
		authService: authService,
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
	if validUser.Username != "" {
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

func (uc *userController) LoginUser(c *gin.Context) {
	var userLoginDTO DTO.LoginUserDTO

	if err := c.ShouldBindJSON(&userLoginDTO); err != nil {
		res := helper.BuildErrorResponse("정상적인 데이터가 아닙니다.", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	checkuser, err := uc.userService.GetUserByUsername(userLoginDTO.Username)
	if err != nil {
		res := helper.BuildErrorResponse("사용자를 찾을 수 없습니다.", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err = uc.userService.ComparePasswords(checkuser.Password, userLoginDTO.Password)
	if err != nil {
		res := helper.BuildErrorResponse("사용자 비밀번호가 일치하지 않습니다.", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	token, err := uc.authService.CreateUserLoginJWT(userLoginDTO)
	if err != nil {
		res := helper.BuildErrorResponse("사용자 토큰을 생성 할 수 없습니다.", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "로그인 성공", token)
	c.JSON(http.StatusOK, res)
}
