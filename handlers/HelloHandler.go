package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func WelcomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello from Gin!")
}
