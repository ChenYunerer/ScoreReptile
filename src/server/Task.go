package server

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/model/http"
	"github.com/gin-gonic/gin"
)

func getTopTaskList(c *gin.Context) {
	var taskInfoList []model.ReptileTaskInfo
	db.Engine.Desc(model.CreateTime).Find(&taskInfoList)
	// 任务分组
	var topTaskList []*http.ReptileTaskWithChild
	for _, item := range taskInfoList {
		if item.Parent_task_id == "" {
			topTaskInfo := genReptileTaskWithChild(item)
			genTaskDetail(topTaskInfo, taskInfoList)
			topTaskList = append(topTaskList, topTaskInfo)
		}
	}
	c.JSON(200, http.GenSuccessResponseWithData(topTaskList))
}

func genTaskDetail(taskInfo *http.ReptileTaskWithChild, taskInfoList []model.ReptileTaskInfo) {
	subTaskInfoList := getTaskInfoListByParentTaskId(taskInfoList, taskInfo.Task_id)
	if len(subTaskInfoList) == 0 {
		return
	}
	taskInfo.SubTaskList = genReptileTaskWithChildList(subTaskInfoList)
	for _, item := range taskInfo.SubTaskList {
		genTaskDetail(item, taskInfoList)
	}
}

func getTaskInfoListByParentTaskId(taskInfo []model.ReptileTaskInfo, parentTaskId string) []model.ReptileTaskInfo {
	taskInfoList := make([]model.ReptileTaskInfo, 0)
	for _, item := range taskInfo {
		if item.Parent_task_id == parentTaskId {
			taskInfoList = append(taskInfoList, item)
		}
	}
	return taskInfoList
}

func genReptileTaskWithChild(taskInfo model.ReptileTaskInfo) *http.ReptileTaskWithChild {
	return &http.ReptileTaskWithChild{
		ReptileTaskInfo: taskInfo,
	}
}

func genReptileTaskWithChildList(taskInfoList []model.ReptileTaskInfo) []*http.ReptileTaskWithChild {
	resultList := make([]*http.ReptileTaskWithChild, 0)
	for _, item := range taskInfoList {
		resultList = append(resultList, &http.ReptileTaskWithChild{
			ReptileTaskInfo: item,
		})
	}
	return resultList
}
