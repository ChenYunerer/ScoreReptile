package server

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/model/http"
	"github.com/gin-gonic/gin"
)

func getTaskScoreList(c *gin.Context) {
	topTaskId := c.Query("topTaskId")
	if topTaskId == "" {
		c.JSON(200, http.GenErrorResponse("topTaskId is empty"))
		return
	}

	var scoreBaseInfoList []model.ScoreBaseInfo
	db.Engine.Where(model.TopTaskId+"= ?", topTaskId).Desc(model.ScoreUploadTime).Find(&scoreBaseInfoList)
	c.JSON(200, http.GenSuccessResponseWithData(scoreBaseInfoList))
}
