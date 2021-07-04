package model

type ScorePictureInfo struct {
	Id                int    `json:"id"`                                           //ID
	ScoreId           int    `xorm:"score_id" json:"scoreId"`                      //曲谱id
	ScoreName         string `xorm:"score_name" json:"scoreName"`                  //曲谱名称
	ScoreHref         string `xorm:"score_href" json:"scoreHref"`                  //曲谱地址
	ScorePictureIndex int    `xorm:"score_picture_index" json:"scorePictureIndex"` //曲谱图片index
	ScorePictureHref  string `xorm:"score_picture_href" json:"scorePictureHref"`   //曲谱图片href
	TopTaskId         string `xorm:"top_task_id" json:"topTaskId"`                 //所属顶层父任务ID
	TaskId            string `xorm:"task_id" json:"taskId"`                        //所属任务ID
}

var ScorePictureIndex = "score_picture_index"
