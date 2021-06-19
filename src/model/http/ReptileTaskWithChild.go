package http

import "ScoreReptile/src/model"

type ReptileTaskWithChild struct {
	model.ReptileTaskInfo
	SubTaskList []*ReptileTaskWithChild `json:"subTaskList"`
}
