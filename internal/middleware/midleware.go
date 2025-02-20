package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinCheckLogin(next gin.HandlerFunc) gin.HandlerFunc {
	f := func(c *gin.Context) {
		if GinCheckRequesCookie(c) {
			next(c)
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/v2/login")
			c.Abort()
		}
	}
	return f
}

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("\t Check Auth")
		if GinCheckRequesCookie(c) {
			c.Next()
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/v2/login")
			c.Abort()
		}
	}
}
