package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/net"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"time"
)

/**
 * 获取曲谱列表
 */

const BaseUrl = "http://www.qupu123.com/"

var categoryDefinitionList []CategoryDefinition

type CategoryDefinition struct {
	categoryName         string
	url                  string
	tempScoreReptileFunc func(url, category string, scoreListTempList *[]*model.ScoreListTemp)
}

func init() {
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "制谱园地",
		url:                  BaseUrl + "jipu",
		tempScoreReptileFunc: tempScoreReptileListType1,
	})
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "原创专栏",
		url:                  BaseUrl + "yuanchuang",
		tempScoreReptileFunc: tempScoreReptileListType1,
	})
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "器乐",
		url:                  BaseUrl + "qiyue",
		tempScoreReptileFunc: tempScoreReptileListType2,
	})
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "戏曲",
		url:                  BaseUrl + "xiqu",
		tempScoreReptileFunc: tempScoreReptileListType2,
	})
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "谱友园地",
		url:                  BaseUrl + "puyou",
		tempScoreReptileFunc: tempScoreReptileListType2,
	})
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "民歌",
		url:                  BaseUrl + "minge",
		tempScoreReptileFunc: tempScoreReptileListType3,
	})
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "美声",
		url:                  BaseUrl + "meisheng",
		tempScoreReptileFunc: tempScoreReptileListType3,
	})
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "通俗",
		url:                  BaseUrl + "tongsu",
		tempScoreReptileFunc: tempScoreReptileListType3,
	})
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "外国",
		url:                  BaseUrl + "waiguo",
		tempScoreReptileFunc: tempScoreReptileListType3,
	})
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "少儿",
		url:                  BaseUrl + "shaoer",
		tempScoreReptileFunc: tempScoreReptileListType3,
	})
	categoryDefinitionList = append(categoryDefinitionList, CategoryDefinition{
		categoryName:         "合唱",
		url:                  BaseUrl + "hechang",
		tempScoreReptileFunc: tempScoreReptileListType3,
	})
}

//jipu（制谱园地） yuanchuang(原创专栏) qiyue（器乐）xiqu（戏曲）puyou（谱友园地）
//声乐 minge（民歌）meisheng（美声）tongsu（通俗）waiguo（外国）shaoer（少儿）hechang（合唱）
//由于html排版不同需要区别处理每个大类的数据
//jipu（制谱园地） yuanchuang(原创专栏) 使用tempScoreReptileListType1
//qiyue（器乐）xiqu（戏曲）puyou(谱友园地) 使用tempScoreReptileListType2
//minge（民歌）meisheng（美声）tongsu（通俗）waiguo（外国）shaoer（少儿）hechang（合唱） 使用tempScoreReptileListType3
func startProcessListTemp(parentTaskInfo model.ReptileTaskInfo) model.ReptileTaskInfo {
	// 创建并插入任务
	taskInfo := model.CreateBasicTaskInfo("各分类列表任务")
	taskInfo.Top_task_id = parentTaskInfo.Task_id
	taskInfo.Parent_task_id = parentTaskInfo.Task_id
	db.Engine.InsertOne(taskInfo)

	//开始执行各个子任务流程
	taskWrapperList := make([]model.ReptileTaskWrapper, 0)

	wg := waitGroup
	for _, item := range categoryDefinitionList {
		wg.Add(1)
		go func(definition CategoryDefinition) {
			defer wg.Done()
			taskWrapper := processChildTask(definition.categoryName+"-ListTempTask", *taskInfo, func(scoreListTempList *[]*model.ScoreListTemp) {
				definition.tempScoreReptileFunc(definition.url, definition.categoryName, scoreListTempList)
			})
			taskWrapperList = append(taskWrapperList, taskWrapper)
		}(item)
	}
	wg.Wait()

	//各个子任务流程结束 插入子任务信息插入子任务数据
	//过滤重复数据
	set := make(map[string]*model.ScoreListTemp)
	for _, taskWrapper := range taskWrapperList {
		for _, scoreListTemp := range taskWrapper.ScoreListTempList {
			_, exist := set[scoreListTemp.ScoreHref]
			if !exist {
				_, err := db.Engine.InsertOne(scoreListTemp)
				if err != nil {
					log.Println("InsertOne scoreListTemp err:", err)
				}
				set[scoreListTemp.ScoreHref] = scoreListTemp
			}
		}
		_, err := db.Engine.InsertOne(taskWrapper.ReptileTaskInfo)
		if err != nil {
			log.Println("InsertOne ReptileTaskInfo err:", err)
		}
	}
	//更新父任务信息
	endTime := time.Now()
	taskInfo.Task_status = 2
	taskInfo.Task_end_time = endTime
	taskInfo.Update_time = endTime
	taskInfo.Task_time_consume = taskInfo.Task_end_time.Sub(taskInfo.Task_start_time).Seconds()
	taskProcessDataNum := 0
	for _, item := range taskWrapperList {
		taskProcessDataNum = taskProcessDataNum + item.ReptileTaskInfo.Task_process_data_num
	}
	taskInfo.Task_process_data_num = taskProcessDataNum
	db.Engine.Update(taskInfo, &model.ReptileTaskInfo{
		Task_id: taskInfo.Task_id,
	})
	return *taskInfo
}

func processChildTask(taskName string, parentTaskInfo model.ReptileTaskInfo, tempScoreReptileFunc func(scoreListTempList *[]*model.ScoreListTemp)) model.ReptileTaskWrapper {
	taskInfo := model.CreateBasicTaskInfo(taskName)
	taskInfo.Top_task_id = parentTaskInfo.Top_task_id
	taskInfo.Parent_task_id = parentTaskInfo.Task_id
	scoreListTempList := make([]*model.ScoreListTemp, 0)
	tempScoreReptileFunc(&scoreListTempList)
	for i := 0; i < len(scoreListTempList); i++ {
		item := scoreListTempList[i]
		item.TopTaskId = parentTaskInfo.Top_task_id
		item.TaskId = taskInfo.Task_id
	}
	taskInfo.Task_process_data_num = len(scoreListTempList)
	taskInfo.Task_status = 2
	endTime := time.Now()
	timeConsume := endTime.Sub(taskInfo.Task_start_time)
	taskInfo.Task_end_time = endTime
	taskInfo.Task_time_consume = timeConsume.Seconds()
	return model.CreateReptileTaskWrapper(*taskInfo, scoreListTempList)
}

func tempScoreReptileListType1(url, category string, scoreListTempList *[]*model.ScoreListTemp) {
	reader, err := net.GetRequestForReader(url)
	if err != nil {
		log.Println(err)
		return
	}
	document, _ := goquery.NewDocumentFromReader(reader)
	selection := document.Find("tbody tr")
	selections := make([]*goquery.Selection, 0)
	selection.Each(func(index int, s *goquery.Selection) {
		selections = append(selections, s)
	})
	for _, s := range selections {
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
			log.Println(name, href, uploader, author, singer, uploadTime)
			exist := db.IsScoreListTempExist(scoreListTemp)
			if exist {
				log.Println("数据已存在 结束！！！", scoreListTemp)
				return
			}
			*scoreListTempList = append(*scoreListTempList, &scoreListTemp)
		}
	}
	//confirm it have next page
	haveNextPage := strings.Contains(document.Find(".pageHtml").Text(), "下一页")
	if haveNextPage {
		document.Find(".pageHtml").Children().Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "下一页" {
				nextPageHref, _ := selection.Attr("href")
				tempScoreReptileListType1(BaseUrl+nextPageHref, category, scoreListTempList)
			}
		})
	} else {
		log.Println("已经到最后一页了")
	}
}

func tempScoreReptileListType2(url, category string, scoreListTempList *[]*model.ScoreListTemp) {
	reader, err := net.GetRequestForReader(url)
	if err != nil {
		log.Println(err)
		return
	}
	document, _ := goquery.NewDocumentFromReader(reader)
	selection := document.Find("tbody tr")
	selections := make([]*goquery.Selection, 0)
	selection.Each(func(index int, s *goquery.Selection) {
		selections = append(selections, s)
	})
	for _, s := range selections {
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
			log.Println(name, href, uploader, author, singer, uploadTime)
			exist := db.IsScoreListTempExist(scoreListTemp)
			if exist {
				log.Println("数据已存在 结束！！！", scoreListTemp)
				return
			}
			*scoreListTempList = append(*scoreListTempList, &scoreListTemp)
		}
	}
	//confirm it have next page
	haveNextPage := strings.Contains(document.Find(".pageHtml").Text(), "下一页")
	if haveNextPage {
		document.Find(".pageHtml").Children().Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "下一页" {
				nextPageHref, _ := selection.Attr("href")
				tempScoreReptileListType2(BaseUrl+nextPageHref, category, scoreListTempList)
			}
		})
	} else {
		log.Println("已经到最后一页了")
	}
}

func tempScoreReptileListType3(url, category string, scoreListTempList *[]*model.ScoreListTemp) {
	reader, err := net.GetRequestForReader(url)
	if err != nil {
		log.Println(err)
		return
	}
	document, _ := goquery.NewDocumentFromReader(reader)
	selection := document.Find("tbody tr")
	selections := make([]*goquery.Selection, 0)
	selection.Each(func(index int, s *goquery.Selection) {
		selections = append(selections, s)
	})
	for _, s := range selections {
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
			log.Println(name, href, uploader, author, singer, uploadTime)
			exist := db.IsScoreListTempExist(scoreListTemp)
			if exist {
				log.Println("数据已存在 结束！！！", scoreListTemp)
				return
			}
			*scoreListTempList = append(*scoreListTempList, &scoreListTemp)
		}
	}
	//confirm it have next page
	haveNextPage := strings.Contains(document.Find(".pageHtml").Text(), "下一页")
	if haveNextPage {
		document.Find(".pageHtml").Children().Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "下一页" {
				nextPageHref, _ := selection.Attr("href")
				tempScoreReptileListType3(BaseUrl+nextPageHref, category, scoreListTempList)
			}
		})
	} else {
		log.Println("已经到最后一页了")
	}
}
