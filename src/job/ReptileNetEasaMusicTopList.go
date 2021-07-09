package job

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/util"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"time"
)

var netEaseTopListBaseUrl = "https://music.163.com/discover/toplist"
var topListList = make([]TopList, 0)

type TopList struct {
	topListId   int
	topListName string
}

func init() {
	topListList = append(topListList, TopList{
		topListId:   3778678,
		topListName: "网易云热歌榜",
	})
}

func startReptileNetEaseMusicTopList(parentTaskInfo model.ReptileTaskInfo) model.ReptileTaskInfo {
	// 创建并插入任务
	taskInfo := model.CreateBasicTaskInfo("榜单任务", parentTaskInfo.Task_type)
	taskInfo.Top_task_id = parentTaskInfo.Task_id
	taskInfo.Parent_task_id = parentTaskInfo.Task_id
	db.Engine.InsertOne(taskInfo)

	wg := waitGroup
	allTopListInfoList := make([]model.NeteaseMusicTopListDetailInfo, 0)
	childTaskList := make([]model.ReptileTaskInfo, 0)
	for _, item := range topListList {
		wg.Add(1)
		go func(topList TopList) {
			defer wg.Done()
			childTaskInfo := model.CreateBasicTaskInfo(topList.topListName+"任务", parentTaskInfo.Task_type)
			childTaskInfo.Top_task_id = taskInfo.Top_task_id
			childTaskInfo.Parent_task_id = taskInfo.Task_id
			topListInfoList := doReptileNetEaseMusicTopList(netEaseTopListBaseUrl, topList.topListId, topList.topListName, *childTaskInfo)
			childTaskInfo.Task_process_data_num = len(topListInfoList)
			childTaskInfo.Task_status = 2
			endTime := time.Now()
			timeConsume := endTime.Sub(taskInfo.Task_start_time)
			childTaskInfo.Task_end_time = endTime
			childTaskInfo.Task_time_consume = timeConsume.Seconds()
			childTaskList = append(childTaskList, *childTaskInfo)
			allTopListInfoList = append(allTopListInfoList, topListInfoList...)
		}(item)
	}
	wg.Wait()
	// 子任务信息入库
	db.Engine.Insert(childTaskList)
	// 榜单数据入库
	_, err := db.Engine.Insert(allTopListInfoList)
	if err != nil {
		log.Println(err)
	}
	// 更新父任务信息
	endTime := time.Now()
	taskInfo.Task_status = 2
	taskInfo.Task_end_time = endTime
	taskInfo.Update_time = endTime
	taskInfo.Task_time_consume = taskInfo.Task_end_time.Sub(taskInfo.Task_start_time).Seconds()
	taskProcessDataNum := 0
	for _, item := range childTaskList {
		taskProcessDataNum = taskProcessDataNum + item.Task_process_data_num
	}
	taskInfo.Task_process_data_num = taskProcessDataNum
	db.Engine.Update(taskInfo, &model.ReptileTaskInfo{
		Task_id: taskInfo.Task_id,
	})
	return *taskInfo
}

func doReptileNetEaseMusicTopList(baseUrl string, topListId int, topListName string, parentTaskInfo model.ReptileTaskInfo) []model.NeteaseMusicTopListDetailInfo {
	nowTime := time.Now()
	topListInfoList := make([]model.NeteaseMusicTopListDetailInfo, 0)
	url := baseUrl + "?id=" + strconv.Itoa(topListId)
	reader, err := util.GetRequestForReader(url)
	if err != nil {
		log.Println(err)
	}
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err)
	}
	// 获取更新时间
	originTopListUpdateTime := document.Find(".sep, .s-fc3").First().Text()
	topListUpdateTimeStr := strings.ReplaceAll(originTopListUpdateTime, "最近更新：", "")
	topListUpdateTimeStr = strings.ReplaceAll(topListUpdateTimeStr, "日", "")
	splits := strings.Split(topListUpdateTimeStr, "月")
	month, _ := strconv.Atoi(splits[0])
	day, _ := strconv.Atoi(splits[1])
	topListUpdateTime := time.Date(nowTime.Year(), time.Month(month), day, 0, 0, 0, 0, nowTime.Location())
	// 判断日期是否存在
	count, _ := db.Engine.Where(model.TopListId+" = ? and "+model.TopListUpdateTime+" = ?", topListId, topListUpdateTime.String()).Count(&model.NeteaseMusicTopListDetailInfo{})
	if count > 0 {
		return topListInfoList
	}
	// 获取榜单数据
	document.Find("#song-list-pre-cache li").Each(func(i int, selection *goquery.Selection) {
		hrefStr, _ := selection.Find("a").Attr("href")
		songName := selection.Find("a").Text()
		netEaseMusicTopListInfo := model.NeteaseMusicTopListDetailInfo{
			SongName:          songName,
			Sort:              i,
			Href:              hrefStr,
			TopListId:         topListId,
			TopListName:       topListName,
			TopListUpdateTime: topListUpdateTime,
			CreateTime:        nowTime,
			UpdateTime:        nowTime,
			TopTaskId:         parentTaskInfo.Top_task_id,
			TaskId:            parentTaskInfo.Task_id,
		}
		topListInfoList = append(topListInfoList, netEaseMusicTopListInfo)
		log.Println(topListName, i, hrefStr, songName)
	})
	return topListInfoList
}
