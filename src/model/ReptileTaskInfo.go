package model

import (
	"github.com/google/uuid"
	"time"
)

type ReptileTaskInfo struct {
	Task_id               string    `xorm:"task_id" json:"taskId"`
	Top_task_id           string    `xorm:"top_task_id" json:"topTaskId"`
	Parent_task_id        string    `xorm:"parent_task_id" json:"parentTaskId"`
	Task_type             int       `xorm:"task_type" json:"taskType"` //任务类型 1 曲谱123爬虫任务 2 网易云榜单爬虫任务
	Task_name             string    `xorm:"task_name" json:"taskName"`
	Task_status           int       `xorm:"task_status" json:"taskStatus"` //任务状态 -1失败 0初始化 1处理中 2成功
	Task_process_data_num int       `xorm:"task_process_data_num" json:"taskProcessDataNum"`
	Task_start_time       time.Time `xorm:"task_start_time" json:"taskStartTime"`
	Task_end_time         time.Time `xorm:"task_end_time" json:"taskEndTime"`
	Task_time_consume     float64   `xorm:"task_time_consume" json:"taskTimeConsume"`
	Create_time           time.Time `xorm:"create_time" json:"createTime"`
	Update_time           time.Time `xorm:"update_time" json:"updateTime"`
}

func CreateBasicTaskInfo(name string, taskType int) *ReptileTaskInfo {
	startTime := time.Now()
	_uuid, _ := uuid.NewUUID()
	taskInfo := &ReptileTaskInfo{
		Task_id:               _uuid.String(),
		Top_task_id:           _uuid.String(),
		Task_type:             taskType,
		Task_name:             name,
		Task_status:           1,
		Task_process_data_num: 0,
		Task_start_time:       startTime,
		Task_end_time:         startTime,
		Task_time_consume:     0,
		Create_time:           startTime,
		Update_time:           startTime,
	}
	return taskInfo
}

const TaskType = "task_type"
const ParentTaskId = "parent_task_id"
const TaskStatus = "task_status"
const CreateTime = "create_time"
