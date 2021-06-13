package model

type ScoreListTemp struct {
	Id                 int
	ScoreCategory      string `xorm:"score_category"`
	ScoreName          string `xorm:"score_name"`
	ScoreHref          string `xorm:"score_href"`
	ScoreUploader      string `xorm:"score_uploader"`
	ScoreAuthor        string `xorm:"score_author"`
	ScoreSinger        string `xorm:"score_singer"`
	ScoreUploadTime    string `xorm:"score_upload_time"`
	ScoreReptileStatus int    `xorm:"score_reptile_status"` //曲谱爬虫状态 0 未爬 1 已爬
	TopTaskId          string `xorm:"top_task_id"`          //所属顶层父任务ID
	TaskId             string `xorm:"task_id"`              //所属任务ID
}

const TopTaskId = "top_task_id"
const TaskId = "task_id"
