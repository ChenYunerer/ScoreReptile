package server

import (
	"ScoreReptile/src/job"
	"github.com/gin-gonic/gin"
)

func test(c *gin.Context) {
	job.StartNetEaseMusicReptileTask()
	c.JSON(200, "123")
}
