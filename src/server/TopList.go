package server

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/model/http"
	"github.com/gin-gonic/gin"
)

func getTaskTopListDetailList(c *gin.Context) {
	topTaskId := c.Query("topTaskId")
	topListUpdateTime := c.Query("topListUpdateTime")
	var topListDetailList []model.NeteaseMusicTopListDetailInfo
	session := db.Engine.Asc(model.Sort)
	if topTaskId != "" {
		session.Where(model.TopTaskId+"= ?", topTaskId)
	}
	if topListUpdateTime != "" {
		session.Where(model.TopListUpdateTime+"= ?", topListUpdateTime)
	}
	session.Find(&topListDetailList)
	c.JSON(200, http.GenSuccessResponseWithData(topListDetailList))
}

func getTopListAllDate(c *gin.Context) {
	var topListDetailList []model.NeteaseMusicTopListDetailInfo
	db.Engine.Select(model.TopListUpdateTime).GroupBy(model.TopListUpdateTime).Find(&topListDetailList)
	timeList := make([]string, 0)
	for _, item := range topListDetailList {
		timeList = append(timeList, item.TopListUpdateTime.String())
	}
	c.JSON(200, http.GenSuccessResponseWithData(timeList))
}
