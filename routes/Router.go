package routes

import (
	"github.com/gin-gonic/gin"
	"playground/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./assets/ico/favicon.ico")
	})

	r.GET("/", handlers.WelcomeHandler)

	apiGroup := r.Group("/api")
	apiGroup.GET("/", handlers.WelcomeHandler)
	apiGroup.GET("/hello", handlers.HelloHandler)

	return r
}
