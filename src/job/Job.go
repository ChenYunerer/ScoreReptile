package job

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func StartCronJob() {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 0 1 * * ?", func() {
		startReptileTask()
	})
	if err != nil {
		log.Panic("add cron func err: ", err)
	}
	c.Start()
}

func startReptileTask() {
	//创建最顶级任务
	nowTime := time.Now()
	taskName := fmt.Sprintf("%d-%d-%d 爬虫任务", nowTime.Year(), nowTime.Month(), nowTime.Day())
	taskInfo := model.CreateBasicTaskInfo(taskName)
	if _, err := db.Engine.InsertOne(taskInfo); err != nil {
		log.Println(err)
	}
	childTaskInfoList := make([]model.ReptileTaskInfo, 0)

	//执行第一步操作 抓取各个类目的列表数据
	processListTempTaskInfo := startProcessListTemp(*taskInfo)
	childTaskInfoList = append(childTaskInfoList, processListTempTaskInfo)

	//第二步操作 抓取曲谱基本信息
	processBaseInfoTaskInfo := startProcessBaseInfo(*taskInfo)
	childTaskInfoList = append(childTaskInfoList, processBaseInfoTaskInfo)

	//第三步操作 抓取曲谱图片
	processPictureInfoTaskInfo := startProcessPictureInfo(*taskInfo)
	childTaskInfoList = append(childTaskInfoList, processPictureInfoTaskInfo)

	//第四步操作 计算曲谱图片数量
	processScorePictureCountTaskInfo := startProcessScorePictureCount(*taskInfo)
	childTaskInfoList = append(childTaskInfoList, processScorePictureCountTaskInfo)

	//更新任务
	taskInfo.Task_status = 2
	endTime := time.Now()
	taskInfo.Task_end_time = endTime
	timeConsume := endTime.Sub(taskInfo.Task_start_time)
	taskInfo.Task_time_consume = timeConsume.Seconds()
	taskInfo.Update_time = endTime
	taskProcessDataNum := 0
	for _, item := range childTaskInfoList {
		taskProcessDataNum = taskProcessDataNum + item.Task_process_data_num
	}
	taskInfo.Task_process_data_num = taskProcessDataNum
	db.Engine.Update(taskInfo, &model.ReptileTaskInfo{
		Task_id: taskInfo.Task_id,
	})
}
