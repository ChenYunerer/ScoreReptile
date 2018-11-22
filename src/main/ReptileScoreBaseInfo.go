package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/net"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

/**
 * 获取曲谱信息
 */

var scoreBaseInfoChain = make(chan model.ScoreBaseInfo, 2000)

func startProcessBaseInfo() {
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
		//查询该数据是否已经处理过
		exist := db.IsScoreBaseInfoExist(href)
		if exist {
			log.Println("数据已处理 跳过...")
			continue
		}
		//获取HTML
		reader, err := net.GetRequestForReader(BaseUrl + href)
		if err != nil {
			log.Println(err)
			continue
		}
		var scoreBaseInfo model.ScoreBaseInfo
		//封住已知原始数据
		scoreBaseInfo.ScoreName = s.ScoreName
		scoreBaseInfo.ScoreHref = s.ScoreHref
		scoreBaseInfo.ScoreAuthor = s.ScoreAuthor
		scoreBaseInfo.ScoreCategory = s.ScoreCategory
		scoreBaseInfo.ScoreSinger = s.ScoreSinger
		scoreBaseInfo.ScoreUploader = s.ScoreUploader
		scoreBaseInfo.ScoreUploadTime = s.ScoreUploadTime
		//解析HTML
		document, _ := goquery.NewDocumentFromReader(reader)
		selection := document.Find(".content .content_head")
		fullName := selection.Find("h1").Text()
		scoreBaseInfo.ScoreName = fullName
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
		//获取曲谱Id
		//从href获取id
		id := GetScoreIdFromHref(href)
		if id == 0 {
			log.Println("无法从href中获取id，尝试从onlick事件中获取")
		}
		//从onclick中获取曲谱Id
		onclick, exist := document.Find("#look_all").Attr("onclick")
		if !exist {
			log.Println("无法获取曲谱id")
		} else {
			id, _ = strconv.Atoi(strings.Split(strings.Split(onclick, "','")[0], "'")[1])
		}
		scoreBaseInfo.ScoreId = id
		//获取曲谱封面图
		scoreBaseInfo.ScoreViewCount = getScoreViewCount(id)
		scoreCoverPicture, _ := document.Find(".imageList a").Attr("href")
		scoreBaseInfo.ScoreCoverPicture = scoreCoverPicture

		scoreBaseInfoChain <- scoreBaseInfo
	}
}

//从href中解析id
func GetScoreIdFromHref(href string) int {
	strs := strings.Split(href, "/")
	strs = strings.Split(strs[len(strs)-1], "p")
	if len(strs) != 2 {
		return 0
	}
	strs = strings.Split(strs[1], ".")
	if len(strs) != 2 {
		return 0
	}
	id, _ := strconv.Atoi(strs[0])
	return id
}

//获取曲谱浏览量
func getScoreViewCount(id int) int {
	viewCountStr, _ := net.GetRequest(BaseUrl + "Opern-cnum-id-" + strconv.Itoa(id) + ".html")
	viewCount, _ := strconv.Atoi(viewCountStr)
	return viewCount
}
