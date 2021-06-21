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
		taskGroup.GET("/getTopTaskList", getTopTaskList)
	}
	scoreGroup := r.Group("/admin/api")
	{
		scoreGroup.GET("/getTaskScoreList", getTaskScoreList)
	}
	err := r.Run("0.0.0.0:7002")
	if err != nil {
		log.Panic("http server start err: ", err)
	}
}
