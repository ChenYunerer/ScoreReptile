package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/net"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
)

func startProcessPictureInfo() {
	//scorePictureInfoReptile()
	processScorePictureCount()
}

func scorePictureInfoReptile() {
	scoreBaseInfos, err := db.GetScoreBaseInfo(300000)
	if err != nil {
		log.Panic(err)
	}
	for index, s := range scoreBaseInfos {
		log.Println("data index : ", index, " name ", s.ScoreName, " href ", s.ScoreHref)
		reader, err := net.GetRequestForReader(BaseUrl + "Mobile-view-id-" + strconv.Itoa(s.ScoreId) + ".html")
		if err != nil {
			log.Println(err)
			continue
		}
		document, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			log.Println(err)
			continue
		}
		document.Find(".image_list").Find("a").Each(func(i int, selection *goquery.Selection) {
			pictureHref, _ := selection.Attr("href")
			log.Println(s.ScoreName, " ", i, " ", pictureHref)
			scorePictureInfo := model.ScorePictureInfo{
				ScoreId:           s.ScoreId,
				ScoreName:         s.ScoreName,
				ScoreHref:         s.ScoreHref,
				ScorePictureIndex: i,
				ScorePictureHref:  pictureHref,
			}
			err := db.InsertScorePictureInfo(scorePictureInfo)
			if err != nil {
				log.Println("数据库插入失败", err)
			}
		})
	}
}

func processScorePictureCount() {
	scoreBaseInfos, err := db.GetScoreBaseInfo(300000)
	if err != nil {
		log.Panic(err)
	}
	for index, s := range scoreBaseInfos {
		log.Println("update score picture count index : ", index, " name ", s.ScoreName, " href ", s.ScoreHref)
		count := db.CountScorePictureInfo(s.ScoreHref)
		success := db.UpdateScoreBaseInfoPictureCount(s.ScoreHref, count)
		log.Println("score picture count : ", count)
		log.Println(success)
	}
}
