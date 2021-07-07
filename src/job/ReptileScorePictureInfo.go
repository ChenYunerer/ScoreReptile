package job

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/js"
	"ScoreReptile/src/model"
	"ScoreReptile/src/util"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var picGetThreadNum = runtime.NumCPU() * 2

func startProcessPictureInfo(parentTaskInfo model.ReptileTaskInfo) model.ReptileTaskInfo {
	//生成任务
	taskInfo := model.CreateBasicTaskInfo("抓取曲谱图片任务", parentTaskInfo.Task_type)
	taskInfo.Top_task_id = parentTaskInfo.Top_task_id
	taskInfo.Parent_task_id = parentTaskInfo.Task_id
	_, err := db.Engine.InsertOne(taskInfo)
	if err != nil {
		log.Println("InsertOne taskInfo err: ", err)
	}

	//获取需要处理的数据
	var scoreBaseInfoList []model.ScoreBaseInfo
	err = db.Engine.Where(model.TopTaskId+"= ?", taskInfo.Top_task_id).Find(&scoreBaseInfoList)
	if err != nil {
		log.Println("get scoreBaseInfoList err: ", err)
	}
	log.Println("scoreBaseInfoList count： ", len(scoreBaseInfoList))

	//分割数据
	scoreBaseInfosArray := splitScoreBaseInfoArray(scoreBaseInfoList, picGetThreadNum)
	scorePictureInfoList := make([]model.ScorePictureInfo, 0)

	//多线程处理
	wg := waitGroup
	for _, scoreBaseInfos := range scoreBaseInfosArray {
		wg.Add(1)
		go func(items []model.ScoreBaseInfo) {
			defer wg.Done()
			_scorePictureInfoList := pictureInfoReptile(items, *taskInfo)
			scorePictureInfoList = append(scorePictureInfoList, _scorePictureInfoList...)
		}(scoreBaseInfos)
	}
	wg.Wait()

	//插入数据到数据库
	for _, item := range scorePictureInfoList {
		_, err := db.Engine.InsertOne(item)
		if err != nil {
			log.Println("InsertOne scorePictureInfo err: ", err)
		}
	}

	//更新任务信息
	taskInfo.Task_process_data_num = len(scorePictureInfoList)
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

func pictureInfoReptile(scoreBaseInfos []model.ScoreBaseInfo, taskInfo model.ReptileTaskInfo) []model.ScorePictureInfo {
	scorePictureInfoList := make([]model.ScorePictureInfo, 0)
	for index, s := range scoreBaseInfos {
		url := BaseUrl + "Mobile-view-id-" + strconv.Itoa(s.ScoreId) + ".html"
		log.Println("data-index: ", index, " name: ", s.ScoreName, " href: ", s.ScoreHref, " mobile-url: ", url)
		reader, err := util.GetRequestForReader(url)
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
				TopTaskId:         taskInfo.Top_task_id,
			}
			scorePictureInfoList = append(scorePictureInfoList, scorePictureInfo)
		})
	}
	return scorePictureInfoList
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
