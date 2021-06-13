package model

type ScoreBaseInfo struct {
	Id                int
	ScoreId           int    `xorm:"score_id"`            //曲谱id
	ScoreCategory     string `xorm:"score_category"`      //曲谱类别
	ScoreName         string `xorm:"score_name"`          //曲谱名称
	ScoreHref         string `xorm:"score_href"`          //曲谱地址
	ScoreSinger       string `xorm:"score_singer"`        //曲谱演唱者
	ScoreAuthor       string `xorm:"score_author"`        //曲谱词曲作者
	ScoreWordWriter   string `xorm:"score_word_writer"`   //曲谱词作者
	ScoreSongWriter   string `xorm:"score_song_writer"`   //曲谱曲作者
	ScoreFormat       string `xorm:"score_format"`        //曲谱格式
	ScoreOrigin       string `xorm:"score_origin"`        //曲谱来源
	ScoreUploader     string `xorm:"score_uploader"`      //曲谱上传者
	ScoreUploadTime   string `xorm:"score_upload_time"`   //曲谱上传时间
	ScoreViewCount    int    `xorm:"score_view_count"`    //曲谱浏览量
	ScoreCoverPicture string `xorm:"score_cover_picture"` //曲谱封面图（首张曲谱）
	ScorePictureCount int    `xorm:"score_picture_count"` //曲谱图片数量
	TopTaskId         string `xorm:"top_task_id"`         //所属顶层父任务ID
	TaskId            string `xorm:"task_id"`             //所属任务ID
}
