package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"log"
	"runtime"
	"time"
)

var picCountThreadNum = runtime.NumCPU() * 2

func processScorePictureCount(parentTaskInfo model.ReptileTaskInfo) {
	taskInfo := model.CreateBasicTaskInfo("计算曲谱图片任务")
	taskInfo.Top_task_id = parentTaskInfo.Top_task_id
	taskInfo.Parent_task_id = parentTaskInfo.Task_id
	_, err := db.Engine.InsertOne(taskInfo)
	if err != nil {
		log.Println("InsertOne taskInfo err: ", err)
	}

	//获取需要处理的数据
	var scoreBaseInfoList []model.ScoreBaseInfo
	err = db.Engine.Where(model.TopTaskId+"= ?", taskInfo.Top_task_id).Find(&scoreBaseInfoList)
	if err != nil {
		log.Println("get scoreBaseInfoList err: ", err)
	}
	log.Println("scoreBaseInfoList count： ", len(scoreBaseInfoList))

	//分割数据
	scoreBaseInfosList := splitScoreBaseInfoArray(scoreBaseInfoList, picCountThreadNum)

	//多线程处理
	wg := waitGroup
	for _, items := range scoreBaseInfosList {
		wg.Add(1)
		go func(arr []model.ScoreBaseInfo) {
			defer wg.Done()
			countAndUpdatePicCount(arr)
		}(items)
	}
	wg.Wait()

	//更新任务信息
	taskInfo.Task_process_data_num = len(scoreBaseInfoList)
	endTime := time.Now()
	taskInfo.Task_status = 2
	taskInfo.Task_end_time = endTime
	taskInfo.Update_time = endTime
	taskInfo.Task_time_consume = taskInfo.Task_end_time.Sub(taskInfo.Task_start_time).Seconds()
	db.Engine.Update(taskInfo, &model.ReptileTaskInfo{
		Task_id: taskInfo.Task_id,
	})
}

func countAndUpdatePicCount(arr []model.ScoreBaseInfo) {
	for index, s := range arr {
		log.Println("update score picture count index : ", index, " name ", s.ScoreName, " href ", s.ScoreHref)
		count := db.CountScorePictureInfo(s.ScoreHref)
		if count == 0 {
			continue
		}
		success := db.UpdateScoreBaseInfoPictureCount(s.ScoreHref, count)
		log.Println("score picture count : ", count, " update result: ", success)
	}
}
