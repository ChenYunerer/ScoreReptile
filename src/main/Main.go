package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/model/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func main() {
	startCronJob()
	startHttpServer()
}

func startCronJob() {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 0 1 * * ?", func() {
		startReptileTask()
	})
	if err != nil {
		log.Panic("add cron func err: ", err)
	}
	c.Start()
}

func startHttpServer() {
	r := gin.Default()
	r.GET("/a", test2)
	r.GET("/b", test1)
	r.GET("/c", test2)
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		log.Panic("http server start err: ", err)
	}
}

func startReptileTask() {
	//创建最顶级任务
	nowTime := time.Now()
	taskName := fmt.Sprintf("%d-%d-%d 爬虫任务", nowTime.Year(), nowTime.Month(), nowTime.Day())
	taskInfo := model.CreateBasicTaskInfo(taskName)
	if _, err := db.Engine.InsertOne(taskInfo); err != nil {
		log.Println(err)
	}

	//执行第一步操作 抓取各个类目的列表数据
	startProcessListTemp(*taskInfo)

	//第二步操作 抓取曲谱基本信息
	startProcessBaseInfo(*taskInfo)

	//第三步操作 抓取曲谱图片
	startProcessPictureInfo(*taskInfo)

	//第四步操作 计算曲谱图片数量
	processScorePictureCount(*taskInfo)

	//更新任务
	taskInfo.Task_status = 2
	endTime := time.Now()
	taskInfo.Task_end_time = endTime
	timeConsume := endTime.Sub(taskInfo.Task_start_time)
	taskInfo.Task_time_consume = timeConsume.Seconds()
	taskInfo.Update_time = endTime
	db.Engine.Update(taskInfo, &model.ReptileTaskInfo{
		Task_id: taskInfo.Task_id,
	})
}

func test1(c *gin.Context) {
	taskId := c.Query("taskId")
	if taskId == "" {
		c.JSON(200, http.GenErrorResponse("taskId不可为空"))
		return
	}
	var taskInfo model.ReptileTaskInfo
	db.Engine.Where(model.TaskId+"= ?", taskId).Get(&taskInfo)
	//第二部操作 抓取曲谱基本信息
	startProcessBaseInfo(taskInfo)
	c.JSON(200, http.GenSuccessResponse())
}

func test2(c *gin.Context) {
	taskId := c.Query("taskId")
	if taskId == "" {
		c.JSON(200, http.GenErrorResponse("taskId不可为空"))
		return
	}
	var taskInfo model.ReptileTaskInfo
	db.Engine.Where(model.TaskId+"= ?", taskId).Get(&taskInfo)
	//第二部操作 抓取曲谱基本信息
	startProcessBaseInfo(taskInfo)
	c.JSON(200, http.GenSuccessResponse())
}
