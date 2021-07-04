package http

import "ScoreReptile/src/model"

type ScoreDetailInfo struct {
	model.ScoreBaseInfo
	PicInfoList []model.ScorePictureInfo `json:"picInfoList"` //曲谱图片信息
}
