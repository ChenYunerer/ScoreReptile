package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/net"
	"github.com/PuerkitoBio/goquery"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

/**
 * 获取曲谱信息
 */

func startProcessBaseInfo(parentTaskInfo model.ReptileTaskInfo) model.ReptileTaskInfo {
	//生成一个任务
	taskInfo := model.CreateBasicTaskInfo("曲谱基本数据抓取任务")
	taskInfo.Top_task_id = parentTaskInfo.Task_id
	taskInfo.Parent_task_id = parentTaskInfo.Task_id
	db.Engine.InsertOne(taskInfo)

	//获取任务需要处理的数据
	var scoreListTemps []model.ScoreListTemp
	err := db.Engine.Where(model.TopTaskId+"= ?", parentTaskInfo.Task_id).Find(&scoreListTemps)
	if err != nil {
		log.Println("get scoreListTemps err: ", err)
	}
	scoreListTempCount := len(scoreListTemps)
	log.Println("scoreListTempCount： ", scoreListTempCount)

	//将数据分批多线程处理
	threadCount := runtime.NumCPU() * 2
	scoreBaseInfoList := make([]model.ScoreBaseInfo, 0)
	scoreListTempsArray := splitScoreListTempArray(scoreListTemps, threadCount)
	wg := waitGroup
	for _, scoreListTemps := range scoreListTempsArray {
		wg.Add(1)
		go func(scoreListTemps []model.ScoreListTemp) {
			defer wg.Done()
			_scoreBaseInfoList := baseInfoReptile(scoreListTemps, *taskInfo)
			scoreBaseInfoList = append(scoreBaseInfoList, _scoreBaseInfoList...)
		}(scoreListTemps)
	}
	wg.Wait()
	//数据入库
	for _, s := range scoreBaseInfoList {
		log.Println("插入数据", s)
		_, err := db.Engine.Insert(s)
		if err != nil {
			log.Println("Insert scoreBaseInfo err: ", err)
		}
	}

	//更新任务信息
	taskInfo.Task_process_data_num = len(scoreBaseInfoList)
	endTime := time.Now()
	taskInfo.Task_status = 2
	taskInfo.Task_end_time = endTime
	taskInfo.Update_time = endTime
	taskInfo.Task_time_consume = taskInfo.Task_end_time.Sub(taskInfo.Task_start_time).Seconds()
	db.Engine.Update(taskInfo, &model.ReptileTaskInfo{
		Task_id: taskInfo.Task_id,
	})
	return *taskInfo
}

func baseInfoReptile(scoreListTemps []model.ScoreListTemp, taskInfo model.ReptileTaskInfo) []model.ScoreBaseInfo {
	scoreBaseInfoList := make([]model.ScoreBaseInfo, 0)
	for _, s := range scoreListTemps {
		href := s.ScoreHref
		//查询该数据是否已经处理过
		exist := db.IsScoreBaseInfoExist(href)
		if exist {
			log.Println("数据已处理 跳过...")
			continue
		}
		//获取HTML
		log.Println(s)
		reader, err := net.GetRequestForReader(BaseUrl + href)
		if err != nil {
			log.Println(err)
			continue
		}
		var scoreBaseInfo model.ScoreBaseInfo
		//封装已知原始数据
		scoreBaseInfo.TopTaskId = taskInfo.Top_task_id
		scoreBaseInfo.TaskId = taskInfo.Task_id
		scoreBaseInfo.ScoreName = s.ScoreName
		scoreBaseInfo.ScoreHref = s.ScoreHref
		scoreBaseInfo.ScoreAuthor = s.ScoreAuthor
		scoreBaseInfo.ScoreCategory = s.ScoreCategory
		scoreBaseInfo.ScoreSinger = s.ScoreSinger
		scoreBaseInfo.ScoreUploader = s.ScoreUploader
		scoreBaseInfo.ScoreUploadTime = s.ScoreUploadTime
		//解析HTML
		document, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			log.Println(err)
			continue
		}
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

		scoreBaseInfoList = append(scoreBaseInfoList, scoreBaseInfo)
	}
	return scoreBaseInfoList
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

func splitScoreListTempArray(arr []model.ScoreListTemp, num int) [][]model.ScoreListTemp {
	max := len(arr)
	var segmens = make([][]model.ScoreListTemp, 0)
	if max < num {
		return append(segmens, arr)
	}
	quantity := max / num
	end := int(0)
	for i := int(1); i <= num; i++ {
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
