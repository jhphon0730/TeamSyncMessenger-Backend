package routes

import (
	"TeamSyncMessenger-Backend/database"
	"TeamSyncMessenger-Backend/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

var (
	db *sql.DB = database.InitDB()
)

func SetupRouter() (*gin.Engine, *sql.DB) {
	r := gin.Default()

	r.Use(middleware.SetHeader)

	r.GET("/", func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(200, "Good")
	})

	return r, db
}
