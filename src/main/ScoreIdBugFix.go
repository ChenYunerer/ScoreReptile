package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/net"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

func main() {
	hrefs, err := db.GetScoreBaseInfo()
	if err != nil {
		log.Println(err)
		return
	}
	allCount := len(hrefs)
	for index, href := range hrefs {
		log.Println("当前正在处理：第", index, "/", allCount, "项数据 ", href)
		reader, _ := net.GetRequestForReader(BaseUrl + href)
		document, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			log.Println("NewDocumentFromReader 失败 ", err)
			continue
		}
		onclick, exist := document.Find("#look_all").Attr("onclick")
		if !exist {
			log.Println("id解析失败 曲谱可能已经被网站屏蔽")
			continue
		}
		id, _ := strconv.Atoi(strings.Split(strings.Split(onclick, "','")[0], "'")[1])
		log.Println("id:", id)
		if id == 0 {
			continue
		}
		success := db.UpdateScoreBaseInfoId(href, id)
		if success {
			log.Println("更新成功")
		} else {
			log.Println("更新失败")
		}
	}
}
