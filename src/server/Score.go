package server

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/model/http"
	"github.com/gin-gonic/gin"
)

func getScoreGeneralInfo(c *gin.Context) {
	var scoreBaseInfo model.ScoreBaseInfo
	scoreNum, _ := db.Engine.Count(scoreBaseInfo)
	var scorePictureInfo model.ScorePictureInfo
	picNum, _ := db.Engine.Count(scorePictureInfo)
	scoreGeneralInfo := http.ScoreGeneralInfo{
		scoreNum,
		picNum,
	}
	c.JSON(200, http.GenSuccessResponseWithData(scoreGeneralInfo))
}

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

func searchScore(c *gin.Context) {
	searchValue := c.Query("searchValue")
	if searchValue == "" {
		c.JSON(200, http.GenSuccessResponse())
		return
	}

	var scoreBaseInfoList []model.ScoreBaseInfo
	db.Engine.
		Where(model.ScoreName+" like ? or "+model.ScoreSinger+" like ? or "+model.ScoreAuthor+" like ? or "+model.ScoreWordWriter+" like ? or "+model.ScoreSongWriter+" like ? or "+model.ScoreFormat+" like ? or "+model.ScoreUploader+" like ?", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%").
		Desc(model.ScoreUploadTime).
		Limit(30).
		Find(&scoreBaseInfoList)
	c.JSON(200, http.GenSuccessResponseWithData(scoreBaseInfoList))
}

func getScoreDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(200, http.GenErrorResponse("id is empty"))
		return
	}

	var scoreBaseInfoList model.ScoreBaseInfo
	db.Engine.Where(model.Id+"= ?", id).Get(&scoreBaseInfoList)
	var picInfoList []model.ScorePictureInfo
	db.Engine.Where(model.ScoreId+"= ?", scoreBaseInfoList.ScoreId).Asc(model.ScorePictureIndex).Find(&picInfoList)
	scoreDetailInfo := http.ScoreDetailInfo{
		scoreBaseInfoList,
		picInfoList,
	}
	c.JSON(200, http.GenSuccessResponseWithData(scoreDetailInfo))
}
