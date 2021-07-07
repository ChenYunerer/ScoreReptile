package server

import (
	"github.com/gin-gonic/gin"
)

func test(c *gin.Context) {
	c.JSON(200, "123")
}
