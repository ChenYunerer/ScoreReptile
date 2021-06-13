package model

import (
	"github.com/google/uuid"
	"time"
)

type ReptileTaskInfo struct {
	Task_id        string `xorm:"task_id"`
	Top_task_id    string `xorm:"top_task_id"`
	Parent_task_id string `xorm:"parent_task_id"`
	Task_name      string `xorm:"task_name"`
	//任务状态 -1失败 0初始化 1处理中 2成功
	Task_status           int       `xorm:"task_status"`
	Task_process_data_num int       `xorm:"task_process_data_num"`
	Task_start_time       time.Time `xorm:"task_start_time"`
	Task_end_time         time.Time `xorm:"task_end_time"`
	Task_time_consume     float64   `xorm:"task_time_consume"`
	Create_time           time.Time `xorm:"create_time"`
	Update_time           time.Time `xorm:"update_time"`
}

func CreateBasicTaskInfo(name string) *ReptileTaskInfo {
	startTime := time.Now()
	_uuid, _ := uuid.NewUUID()
	taskInfo := &ReptileTaskInfo{
		Task_id:               _uuid.String(),
		Top_task_id:           _uuid.String(),
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
