package model

type ScoreBaseInfo struct {
	Id                int    `json:"id"`
	ScoreId           int    `xorm:"score_id" json:"scoreId"`                      //曲谱id
	ScoreCategory     string `xorm:"score_category" json:"scoreCategory"`          //曲谱类别
	ScoreName         string `xorm:"score_name" json:"scoreName"`                  //曲谱名称
	ScoreHref         string `xorm:"score_href" json:"scoreHref"`                  //曲谱地址
	ScoreSinger       string `xorm:"score_singer" json:"scoreSinger"`              //曲谱演唱者
	ScoreAuthor       string `xorm:"score_author" json:"scoreAuthor"`              //曲谱词曲作者
	ScoreWordWriter   string `xorm:"score_word_writer" json:"scoreWordWriter"`     //曲谱词作者
	ScoreSongWriter   string `xorm:"score_song_writer" json:"scoreSongWriter"`     //曲谱曲作者
	ScoreFormat       string `xorm:"score_format" json:"scoreFormat"`              //曲谱格式
	ScoreOrigin       string `xorm:"score_origin" json:"scoreOrigin"`              //曲谱来源
	ScoreUploader     string `xorm:"score_uploader" json:"scoreUploader"`          //曲谱上传者
	ScoreUploadTime   string `xorm:"score_upload_time" json:"scoreUploadTime"`     //曲谱上传时间
	ScoreViewCount    int    `xorm:"score_view_count" json:"scoreViewCount"`       //曲谱浏览量
	ScoreCoverPicture string `xorm:"score_cover_picture" json:"scoreCoverPicture"` //曲谱封面图（首张曲谱）
	ScorePictureCount int    `xorm:"score_picture_count" json:"scorePictureCount"` //曲谱图片数量
	TopTaskId         string `xorm:"top_task_id" json:"topTaskId"`                 //所属顶层父任务ID
	TaskId            string `xorm:"task_id" json:"taskId"`                        //所属任务ID
}

var Id = "id"
var ScoreId = "score_id"
var ScoreUploadTime = "score_upload_time"
var ScoreName = "score_name"
var ScoreSinger = "score_singer"
var ScoreAuthor = "score_author"
var ScoreWordWriter = "score_word_writer"
var ScoreSongWriter = "score_song_writer"
var ScoreFormat = "score_format"
var ScoreUploader = "score_uploader"
