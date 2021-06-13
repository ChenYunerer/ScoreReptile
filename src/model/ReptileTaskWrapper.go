package model

type ReptileTaskWrapper struct {
	ReptileTaskInfo   ReptileTaskInfo
	ScoreListTempList []*ScoreListTemp
}

func CreateReptileTaskWrapper(reptileTaskInfo ReptileTaskInfo, scoreListTempList []*ScoreListTemp) ReptileTaskWrapper {
	return ReptileTaskWrapper{
		ReptileTaskInfo:   reptileTaskInfo,
		ScoreListTempList: scoreListTempList,
	}
}
