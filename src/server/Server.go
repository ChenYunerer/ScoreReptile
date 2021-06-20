package server

import (
	"github.com/gin-gonic/gin"
	"log"
)

func StartHttpServer() {
	r := gin.Default()
	r.BasePath()
	rg := r.Group("/admin/api")
	{
		rg.GET("/getTopTaskList", getTopTaskList)
	}
	err := r.Run("0.0.0.0:7002")
	if err != nil {
		log.Panic("http server start err: ", err)
	}
}
