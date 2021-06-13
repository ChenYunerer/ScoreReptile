package model

type ScorePictureInfo struct {
	Id                int    //ID
	ScoreId           int    `xorm:"score_id"`            //曲谱id
	ScoreName         string `xorm:"score_name"`          //曲谱名称
	ScoreHref         string `xorm:"score_href"`          //曲谱地址
	ScorePictureIndex int    `xorm:"score_picture_index"` //曲谱图片index
	ScorePictureHref  string `xorm:"score_picture_href"`  //曲谱图片href
	TopTaskId         string `xorm:"top_task_id"`         //所属顶层父任务ID
	TaskId            string `xorm:"task_id"`             //所属任务ID
}
