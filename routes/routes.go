package routes

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(200, "Good")
	})

	return r
}
