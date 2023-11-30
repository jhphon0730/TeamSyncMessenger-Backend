package middleware

import (
	"TeamSyncMessenger-Backend/helper"
	"TeamSyncMessenger-Backend/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

func SetHeader(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0")
	c.Header("Last-Modified", time.Now().String())
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "-1")

	// Pass on to the next-in-chain
	c.Next()
}

func TokenAuthMiddleware(c *gin.Context) {
	token, err := c.Request.Cookie("token")

	if err != nil {
		res := helper.BuildErrorResponse("토큰을 찾을 수 없습니다.", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}
	// Get the JWT string from the cookie
	tknStr := token.Value

	if tknStr == "" {
		res := helper.BuildErrorResponse("토큰을 찾을 수 없습니다.", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	// Initialize a new instance of `Claims`
	claims := &model.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	_, err = jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return model.JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			res := helper.BuildErrorResponse("토큰이 만료 되었습니다.", err.Error(), helper.EmptyObj{})
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		res := helper.BuildErrorResponse("토큰을 찾을 수 없습니다.", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	} else {
		c.Next()
	}
}
