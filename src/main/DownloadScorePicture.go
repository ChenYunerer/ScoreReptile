package main

import (
	"ScoreReptile/src/db"
	"ScoreReptile/src/model"
	"ScoreReptile/src/net"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var threadCount = 12
var limit int64 = 10000
var localPath = "score_picture/"

var index = int64(0)
var scorePictureInfoChain = make(chan model.ScorePictureInfo, limit)
var waitGroup = sync.WaitGroup{}

func startDownloadScorePicture() {
	scorePictureInfos, err := db.GetScorePictureInfo(limit)
	if err != nil {
		log.Panic(err)
	}
	for _, scorePictureInfo := range scorePictureInfos {
		log.Println(scorePictureInfo)
		scorePictureInfoChain <- scorePictureInfo
	}
	for i := 0; i < threadCount; i++ {
		waitGroup.Add(1)
		go func(threadName string) {
			for {
				select {
				case scorePictureInfo := <-scorePictureInfoChain:
					atomic.AddInt64(&index, 1)
					log.Print(threadName, " processing index: ", index)
					reader, err := net.GetRequestForReader(BaseUrl + scorePictureInfo.ScorePictureHref)
					if err != nil {
						log.Println(err)
					} else {
						body, err := ioutil.ReadAll(reader)
						if err != nil {
							log.Println(err)
						} else {
							strs := strings.Split(scorePictureInfo.ScorePictureHref, "/")
							fileName := strs[len(strs)-1]
							var path string
							for i := 0; i < len(strs)-1; i++ {
								path = path + strs[i] + "/"
							}
							if err := write2LocalPath(localPath+path, fileName, body); err != nil {
								log.Println(err)
							}
						}
					}
				default:
					if index >= limit {
						waitGroup.Done()
						break
					}
				}
			}
		}("thread " + strconv.Itoa(i))
	}
	waitGroup.Wait()
}

func write2LocalPath(path, fileName string, byteData []byte) error {
	//双斜杠会导致MAC OS系统文件夹无法打开
	path = strings.Replace(path, "//", "/", -1)
	fileName = strings.Replace(fileName, "//", "/", -1)

	log.Println("path and filename", path, fileName)

	exist, err := fileExists(path + fileName)
	if err != nil {
		log.Println("判断文件存在失败")
		return err
	}
	if exist {
		log.Println("文件已经存在 跳过该文件")
		return nil
	}
	err = os.MkdirAll(path, 0655)
	if err != nil {
		log.Println("MkdirAll失败")
		return err
	}
	fileObj, err := os.Create(path + fileName)

	if err != nil {
		log.Println("创建文件失败")
		return err
	}

	defer fileObj.Close()
	if _, err := fileObj.Write(byteData); err != nil {
		log.Println("写入文件失败")
		return err
	}
	return nil
}

func fileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
