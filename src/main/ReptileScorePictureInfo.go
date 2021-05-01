package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/js"
	"ScoreReptile/src/model"
	"ScoreReptile/src/net"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	"log"
	"strconv"
	"strings"
)

func startProcessPictureInfo() {
	scorePictureInfoReptile()
	processScorePictureCount()
}

func scorePictureInfoReptile() {
	threadCount := 6
	wg := waitGroup
	scoreBaseInfos, err := db.GetUnCountPicScoreBaseInfo()
	if err != nil {
		log.Panic(err)
	}
	scoreBaseInfosArray := splitScoreBaseInfoArray(scoreBaseInfos, threadCount)
	for _, scoreBaseInfos := range scoreBaseInfosArray {
		wg.Add(1)
		go func(items []model.ScoreBaseInfo) {
			defer wg.Done()
			pictureInfoReptile(items)
		}(scoreBaseInfos)
	}
	wg.Wait()
}

func pictureInfoReptile(scoreBaseInfos []model.ScoreBaseInfo) {
	for index, s := range scoreBaseInfos {
		url := BaseUrl + "Mobile-view-id-" + strconv.Itoa(s.ScoreId) + ".html"
		log.Println("data-index: ", index, " name: ", s.ScoreName, " href: ", s.ScoreHref, " mobile-url: ", url)
		reader, err := net.GetRequestForReader(url)
		if err != nil {
			log.Println(err)
			continue
		}
		document, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			log.Println(err)
			continue
		}
		var vm *otto.Otto
		document.Find(".image_list").Find("a, script").Each(func(i int, selection *goquery.Selection) {
			var pictureHref string
			if selection.Is("a") {
				pictureHref, _ = selection.Attr("href")
			} else if selection.Is("script") {
				if vm == nil {
					vm = initJSVm(vm, document)
				}
				value, err := vm.Run(selection.Text())
				if err != nil {
					fmt.Println(err)
				}
				pictureHref = value.String()
			}

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
			} else {
				//log.Println("数据库插入成功")
			}
		})
	}
}

func initJSVm(vm *otto.Otto, document *goquery.Document) *otto.Otto {
	log.Println("init vm")
	vm = otto.New()
	_, err := vm.Run(js.JS)
	if err != nil {
		log.Println(err)
	}
	document.Find("script").Each(func(i int, selection *goquery.Selection) {
		if strings.Contains(selection.Text(), "var") {
			_, err := vm.Run(selection.Text())
			if err != nil {
				log.Println(err)
			}
		}
	})
	return vm
}

func processScorePictureCount() {
	scoreBaseInfos, err := db.GetUnCountPicScoreBaseInfo()
	if err != nil {
		log.Panic(err)
	}
	for index, s := range scoreBaseInfos {
		log.Println("update score picture count index : ", index, " name ", s.ScoreName, " href ", s.ScoreHref)
		count := db.CountScorePictureInfo(s.ScoreHref)
		if count == 0 {
			continue
		}
		success := db.UpdateScoreBaseInfoPictureCount(s.ScoreHref, count)
		log.Println("score picture count : ", count)
		log.Println(success)
	}
}

func splitScoreBaseInfoArray(arr []model.ScoreBaseInfo, num int) [][]model.ScoreBaseInfo {
	max := len(arr)
	var segmens = make([][]model.ScoreBaseInfo, 0)
	if max < num {
		return append(segmens, arr)
	}
	quantity := max / num
	end := 0
	for i := 1; i <= num; i++ {
		qu := i * quantity
		if i != num {
			segmens = append(segmens, arr[i-1+end:qu])
		} else {
			segmens = append(segmens, arr[i-1+end:])
		}
		end = qu - i
	}
	return segmens
}
