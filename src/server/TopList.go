package server

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/model/http"
	"github.com/gin-gonic/gin"
)

func getTaskTopListDetailList(c *gin.Context) {
	topTaskId := c.Query("topTaskId")
	if topTaskId == "" {
		c.JSON(200, http.GenErrorResponse("topTaskId is empty"))
		return
	}
	var topListDetailList []model.NeteaseMusicTopListDetailInfo
	db.Engine.Where(model.TopTaskId+"= ?", topTaskId).Asc(model.Sort).Find(&topListDetailList)
	c.JSON(200, http.GenSuccessResponseWithData(topListDetailList))
}
