package job

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/elastic"
	"encoding/json"
	"log"
)

func Upload2Elastic() {
	scoreBaseInfos, err := db.GetScoreBaseInfo(50000)
	if err != nil {
		log.Panic(err)
	}
	for _, scoreBaseInfo := range scoreBaseInfos {
		docBytes, err := json.Marshal(scoreBaseInfo)
		if err != nil {
			log.Fatal(err)
			continue
		}
		elastic.Index("score_base_info", string(docBytes), scoreBaseInfo.ScoreId)
	}

}
