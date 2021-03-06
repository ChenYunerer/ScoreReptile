package server

import (
	"github.com/gin-gonic/gin"
	"log"
)

func StartHttpServer() {
	r := gin.Default()
	r.BasePath()
	taskGroup := r.Group("/admin/api")
	{
		taskGroup.GET("/getTaskGeneralInfo", getTaskGeneralInfo)
		taskGroup.GET("/getTopTaskList", getTopTaskList)
		taskGroup.GET("/getTaskScoreList", getTaskScoreList)
		taskGroup.GET("/getTaskTopListDetailList", getTaskTopListDetailList)
		taskGroup.GET("/getTopListAllDate", getTopListAllDate)
	}
	scoreGroup := r.Group("/admin/api")
	{
		taskGroup.GET("/getScoreGeneralInfo", getScoreGeneralInfo)
		scoreGroup.GET("/searchScore", searchScore)
		scoreGroup.GET("/getScoreDetail", getScoreDetail)
	}
	debugGroup := r.Group("/admin/api/debug")
	{
		debugGroup.GET("/test", test)
	}
	err := r.Run("0.0.0.0:7002")
	if err != nil {
		log.Panic("http server start err: ", err)
	}
}
