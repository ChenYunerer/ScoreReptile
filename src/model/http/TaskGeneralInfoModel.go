package http

type TaskGeneralInfo struct {
	TaskNum    int64 `json:"taskNum"`
	SuccessNum int64 `json:"successNum"`
	FailNum    int64 `json:"failNum"`
	//LatestTaskInfo model.ReptileTaskInfo `json:"latestTaskInfo"`
}
