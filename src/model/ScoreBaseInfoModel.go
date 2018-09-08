package model

type ScoreBaseInfo struct {
	ScoreId           int    //曲谱id
	ScoreCategory     string //曲谱类别
	ScoreName         string //曲谱名称
	ScoreHref         string //曲谱地址
	ScoreSinger       string //曲谱演唱者
	ScoreAuthor       string //曲谱词曲作者
	ScoreWordWriter   string //曲谱词作者
	ScoreSongWriter   string //曲谱曲作者
	ScoreFormat       string //曲谱格式
	ScoreOrigin       string //曲谱来源
	ScoreUploader     string //曲谱上传者
	ScoreUploadTime   string //曲谱上传时间
	ScoreViewCount    int    //曲谱浏览量
	ScoreCoverPicture string //曲谱封面图（首张曲谱）
	ScorePictureCount int    //曲谱图片数量
}
