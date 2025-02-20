package main

import (
	"onepage/internal/middleware"

	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine, e *Env) {

	router.NoRoute(e.NotFound)

	v2 := router.Group("/v2")
	noAuth := v2.Group("")

	v2.Use(middleware.CheckAuth())

	// роуты имаджборды
	v2.GET("/", e.GetAllBoards)
	v2.GET("/boards", e.GetAllBoards)
	v2.POST("/boards", e.AddBoard)

	v2.GET("/threads/:board", e.GetThreadsFromBoard)
	v2.POST("/threads/:board", e.AddThread)

	v2.GET("/threads/:board/:id", e.GetPostFromThread)
	v2.POST("/threads/:board/:id", e.AddPost)

	noAuth.GET("/reg", e.Register)
	noAuth.GET("/login", e.Login)
	noAuth.POST("/reg", e.PostRegister)
	noAuth.POST("/login", e.PostLogin)

	v2.GET("/test", e.GetTest)
	v2.POST("/test", e.PostTest)

	v2.GET("/flush", e.FlushAllTable)

	noAuth.GET("/t", e.Test)
}
