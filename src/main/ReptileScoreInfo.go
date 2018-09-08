package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/net"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

/**
 * 获取曲谱信息
 */

var scoreBaseInfoChain = make(chan model.ScoreBaseInfo, 2000)

func main() {
	scoreListTempCount := db.CountScoreListTemp()
	log.Println(scoreListTempCount)

	go func() {
		baseInfoReptile()
	}()

	i := 0
	for {
		select {
		case s := <-scoreBaseInfoChain:
			i++
			log.Println("processing data index ", i)
			log.Println("插入数据", s)
			err := db.InsertScoreBaseInfo(s)
			if err != nil {
				log.Println(err)
			} else {
				if db.UpdateScoreListTempStatus(s.ScoreHref) {
					log.Println("更新状态成功")
				} else {
					log.Println("更新状态失败")
				}
			}
		}
	}
}

func baseInfoReptile() {
	scoreListTemps, err := db.GetScoreListTemps()
	if err != nil {
		log.Println(err)
	}
	for _, s := range scoreListTemps {
		href := s.ScoreHref
		exist := db.IsScoreBaseInfoExist(s.ScoreHref)
		if exist {
			log.Println("数据已处理 跳过...")
			continue
		}
		reader, err := net.GetRequestForReader(BaseUrl + href)
		if err != nil {
			log.Println(err)
			continue
		}
		var scoreBaseInfo model.ScoreBaseInfo
		id, err := getScoreIdFromHref(href)
		if err != nil {
			log.Println("id解析失败", err)
		}
		scoreBaseInfo.ScoreId = id
		scoreBaseInfo.ScoreName = s.ScoreName
		scoreBaseInfo.ScoreHref = s.ScoreHref
		scoreBaseInfo.ScoreAuthor = s.ScoreAuthor
		scoreBaseInfo.ScoreCategory = s.ScoreCategory
		scoreBaseInfo.ScoreSinger = s.ScoreSinger
		scoreBaseInfo.ScoreUploader = s.ScoreUploader
		scoreBaseInfo.ScoreUploadTime = s.ScoreUploadTime
		document, _ := goquery.NewDocumentFromReader(reader)
		selection := document.Find(".content .content_head")
		//fullName := selection.Find("h1").Text()
		keys := make([]string, 0)
		selection.Find(".info span").Each(func(i int, selection *goquery.Selection) {
			if selection.Text() != "" {
				keys = append(keys, strings.TrimSpace(selection.Text()))
			}
		})
		for index, key := range keys {
			if (index + 1) >= len(keys) {
				break
			}
			strs := strings.Split(selection.Find(".info").Text(), key)
			strs = strings.Split(strs[1], keys[index+1])
			value := strings.TrimSpace(strs[0])
			switch key {
			case "作词：":
				scoreBaseInfo.ScoreWordWriter = value
			case "作曲：":
				scoreBaseInfo.ScoreSongWriter = value
			case "演唱(奏)：":
				scoreBaseInfo.ScoreSinger = value
			case "格式：":
				scoreBaseInfo.ScoreFormat = value
			case "来源：":
				scoreBaseInfo.ScoreOrigin = value
			case "上传：":
				scoreBaseInfo.ScoreUploader = value
			case "日期：":
				scoreBaseInfo.ScoreUploadTime = value
			}
		}
		//获取曲谱封面图
		scoreBaseInfo.ScoreViewCount = getScoreViewCount(id)
		scoreCoverPicture, _ := document.Find(".imageList a").Attr("href")
		scoreBaseInfo.ScoreCoverPicture = scoreCoverPicture
		scoreBaseInfoChain <- scoreBaseInfo
	}
}

//从href中解析id
func getScoreIdFromHref(href string) (int, error) {
	strs := strings.Split(href, "/p")
	if len(strs) != 2 {
		return 0, errors.New("解析Id失败")
	}
	strs = strings.Split(strs[1], ".")
	if len(strs) != 2 {
		return 0, errors.New("解析Id失败")
	}
	id, _ := strconv.Atoi(strs[0])
	return id, nil
}

//获取曲谱浏览量
func getScoreViewCount(id int) int {
	viewCountStr, _ := net.GetRequest(BaseUrl + "Opern-cnum-id-" + strconv.Itoa(id) + ".html")
	viewCount, _ := strconv.Atoi(viewCountStr)
	return viewCount
}
