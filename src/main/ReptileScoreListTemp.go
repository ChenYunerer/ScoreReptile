package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/net"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

/**
 * 获取曲谱列表
 */

const BaseUrl = "http://www.qupu123.com/"

var scoreListTempChain = make(chan model.ScoreListTemp, 2000)

func startProcessListTemp() {
	//jipu（制谱园地） yuanchuang(原创专栏) qiyue（器乐）xiqu（戏曲）puyou（谱友园地）
	//声乐 minge（民歌）meisheng（美声）tongsu（通俗）waiguo（外国）shaoer（少儿）hechang（合唱）
	//由于html排版不同需要区别处理每个大类的数据
	//jipu（制谱园地） yuanchuang(原创专栏) 使用tempScoreReptileListType1
	//qiyue（器乐）xiqu（戏曲）puyou(谱友园地) 使用tempScoreReptileListType2
	//minge（民歌）meisheng（美声）tongsu（通俗）waiguo（外国）shaoer（少儿）hechang（合唱） 使用tempScoreReptileListType3

	/*go func() {
		//jipu（制谱园地）42790
		tempScoreReptileListType1(BaseUrl+"jipu", "制谱园地")
	}()

	go func() {
		//yuanchuang(原创专栏) 16272 ok
		tempScoreReptileListType1(BaseUrl+"yuanchuang", "原创专栏")
	}()

	go func() {
		//qiyue（器乐）49694 ok
		tempScoreReptileListType2(BaseUrl + "qiyue")
	}()

	go func() {
		//xiqu（戏曲）9493
		tempScoreReptileListType2(BaseUrl + "xiqu")
	}()*/

	go func() {
		//puyou（谱友园地）12189
		tempScoreReptileListType2(BaseUrl + "puyou")
	}()

	/*go func() {
		//minge（民歌）66517
		tempScoreReptileListType3(BaseUrl + "minge", "民歌")
	}()*/

	go func() {
		//meisheng（美声）3777
		tempScoreReptileListType3(BaseUrl+"meisheng", "美声")
	}()

	go func() {
		//tongsu（通俗）17484
		tempScoreReptileListType3(BaseUrl+"tongsu", "通俗")
	}()

	go func() {
		//waiguo（外国）6728
		tempScoreReptileListType3(BaseUrl+"waiguo", "外国")
	}()

	go func() {
		//shaoer（少儿）14185
		tempScoreReptileListType3(BaseUrl+"shaoer", "少儿")
	}()

	go func() {
		//hechang（合唱）6906
		tempScoreReptileListType3(BaseUrl+"hechang", "合唱")
	}()

	i := 0
	for {
		select {
		case s := <-scoreListTempChain:
			i++
			log.Println("processing data index ", i)
			exist := db.IsScoreListTempExist(s)
			if !exist {
				log.Println("插入数据", s)
				err := db.InsertScoreListTemp(s)
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Println("数据已存在")
			}
		default:
			log.Println("no data waiting", i)
		}
	}
}

func tempScoreReptileListType1(url, category string) {
	reader, err := net.GetRequestForReader(url)
	if err != nil {
		log.Println(err)
		return
	}
	document, _ := goquery.NewDocumentFromReader(reader)
	selection := document.Find("tbody tr")
	selection.Each(func(index int, s *goquery.Selection) {
		if s.Children().Length() != 1 {
			name := s.Children().Eq(1).Text()
			href, _ := s.Children().Eq(1).Find("a").Attr("href")
			uploader := s.Children().Eq(2).Text()
			author := s.Children().Eq(3).Text()
			singer := s.Children().Eq(4).Text()
			uploadTime := s.Children().Eq(5).Text()
			scoreListTemp := model.ScoreListTemp{
				ScoreCategory:      category,
				ScoreName:          name,
				ScoreHref:          href,
				ScoreUploader:      uploader,
				ScoreAuthor:        author,
				ScoreSinger:        singer,
				ScoreUploadTime:    uploadTime,
				ScoreReptileStatus: 0,
			}
			scoreListTempChain <- scoreListTemp
			log.Println(name, href, uploader, author, singer, uploadTime)
		}
	})
	//confirm it have next page
	haveNextPage := strings.Contains(document.Find(".pageHtml").Text(), "下一页")
	if haveNextPage {
		document.Find(".pageHtml").Children().Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "下一页" {
				nextPageHref, _ := selection.Attr("href")
				tempScoreReptileListType1(BaseUrl+nextPageHref, category)
			}
		})
	}
}

func tempScoreReptileListType2(url string) {
	reader, err := net.GetRequestForReader(url)
	if err != nil {
		log.Println(err)
		return
	}
	document, _ := goquery.NewDocumentFromReader(reader)
	selection := document.Find("tbody tr")
	selection.Each(func(index int, s *goquery.Selection) {
		if s.Children().Length() != 1 {
			name, _ := s.Children().Eq(1).Children().First().Attr("title")
			href, _ := s.Children().Eq(1).Find("a").Attr("href")
			category := s.Children().Eq(2).Text()
			uploader := s.Children().Eq(5).Text()
			author := s.Children().Eq(3).Text()
			singer := s.Children().Eq(4).Text()
			uploadTime := s.Children().Eq(6).Text()
			scoreListTemp := model.ScoreListTemp{
				ScoreCategory:      category,
				ScoreName:          name,
				ScoreHref:          href,
				ScoreUploader:      uploader,
				ScoreAuthor:        author,
				ScoreSinger:        singer,
				ScoreUploadTime:    uploadTime,
				ScoreReptileStatus: 0,
			}
			scoreListTempChain <- scoreListTemp
			log.Println(name, href, uploader, author, singer, uploadTime)
		}
	})
	//confirm it have next page
	haveNextPage := strings.Contains(document.Find(".pageHtml").Text(), "下一页")
	if haveNextPage {
		document.Find(".pageHtml").Children().Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "下一页" {
				nextPageHref, _ := selection.Attr("href")
				tempScoreReptileListType2(BaseUrl + nextPageHref)
			}
		})
	}
}

func tempScoreReptileListType3(url, category string) {
	reader, err := net.GetRequestForReader(url)
	if err != nil {
		log.Println(err)
		return
	}
	document, _ := goquery.NewDocumentFromReader(reader)
	selection := document.Find("tbody tr")
	selection.Each(func(index int, s *goquery.Selection) {
		if s.Children().Length() != 1 {
			name, _ := s.Children().Eq(1).Children().First().Attr("title")
			href, _ := s.Children().Eq(1).Find("a").Attr("href")
			uploader := s.Children().Eq(5).Text()
			author := s.Children().Eq(3).Text()
			singer := s.Children().Eq(4).Text()
			uploadTime := s.Children().Eq(6).Text()
			scoreListTemp := model.ScoreListTemp{
				ScoreCategory:      category,
				ScoreName:          name,
				ScoreHref:          href,
				ScoreUploader:      uploader,
				ScoreAuthor:        author,
				ScoreSinger:        singer,
				ScoreUploadTime:    uploadTime,
				ScoreReptileStatus: 0,
			}
			scoreListTempChain <- scoreListTemp
			log.Println(name, href, uploader, author, singer, uploadTime)
		}
	})
	//confirm it have next page
	haveNextPage := strings.Contains(document.Find(".pageHtml").Text(), "下一页")
	if haveNextPage {
		document.Find(".pageHtml").Children().Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "下一页" {
				nextPageHref, _ := selection.Attr("href")
				tempScoreReptileListType3(BaseUrl+nextPageHref, category)
			}
		})
	}
}
