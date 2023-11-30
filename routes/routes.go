package routes

import (
	"TeamSyncMessenger-Backend/controller"
	"TeamSyncMessenger-Backend/database"
	"TeamSyncMessenger-Backend/middleware"
	"TeamSyncMessenger-Backend/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

var (
	db *sql.DB = database.InitDB()

	userService service.UserService = service.NewUserService(db)

	userController controller.UserController = controller.NewUserController(userService)
)

func SetupRouter() (*gin.Engine, *sql.DB) {
	r := gin.Default()

	r.Use(middleware.SetHeader)

	user_group := r.Group("/api/user")
	{
		user_group.GET("/", userController.GetUsers)
		user_group.POST("/register/", userController.RegisterUser)
	}

	return r, db
}
