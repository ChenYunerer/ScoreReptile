package db

import (
	"ScoreReptile/src/model"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func DeteleScoreListTemps() error {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return err
	}
	defer db.Close()

	res, err := db.Exec("DELETE FROM score_list_temp_tbl")
	if err != nil {
		log.Println("数据库操作失败 : %v", err)
		return err
	}
	rows, errs := res.RowsAffected()
	if rows < 1 {
		log.Println("删除数据库条数小于1原因是 :%v", errs)
		return err
	}
	return nil
}

func InsertScoreListTemp(scoreListTemp model.ScoreListTemp) error {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return err
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO score_list_temp_tbl(score_category, score_name, score_href, score_uploader, score_author, score_singer, score_uploader_time, score_reptile_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", scoreListTemp.ScoreCategory, scoreListTemp.ScoreName, scoreListTemp.ScoreHref, scoreListTemp.ScoreUploader, scoreListTemp.ScoreAuthor, scoreListTemp.ScoreSinger, scoreListTemp.ScoreUploadTime, scoreListTemp.ScoreReptileStatus)
	if err != nil {
		log.Println("数据库添加失败 : %v", err)
		return err
	}
	rows, errs := res.RowsAffected()
	log.Println(res.RowsAffected())
	if rows < 1 {
		log.Println("添加数据库条数小于1原因是 :%v", errs)
		return err
	}

	return nil
}

func IsScoreListTempExist(scoreListTemp model.ScoreListTemp) bool {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return false
	}
	defer db.Close()
	rows := db.QueryRow("select count(*) from score_list_temp_tbl where score_href = ?", scoreListTemp.ScoreHref)
	var count int64
	err = rows.Scan(&count)
	if err != nil {
		log.Println("数据库查询失败 : %v", err)
	}
	return count > 0
}

func GetScoreListTemps() ([]model.ScoreListTemp, error) {
	scoreListTemps := make([]model.ScoreListTemp, 0)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select score_category, score_name, score_href, score_uploader, score_author, score_singer, score_uploader_time, score_reptile_status from score_list_temp_tbl")
	if err != nil {
		log.Println("数据查询库失败 : %v", err)
		return nil, err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		i++
		log.Println("处理数据库数据", i)
		var score_category, score_name, score_href, score_uploader, score_author, score_singer, score_uploader_time string
		var score_reptile_status int
		err = rows.Scan(&score_category, &score_name, &score_href, &score_uploader, &score_author, &score_singer, &score_uploader_time, &score_reptile_status)
		if err != nil {
			log.Println("数据查询对象映射失败 : %v", err)
		} else {
			scoreListTemp := model.ScoreListTemp{
				ScoreCategory:      score_category,
				ScoreName:          score_name,
				ScoreHref:          score_href,
				ScoreUploader:      score_uploader,
				ScoreSinger:        score_singer,
				ScoreAuthor:        score_author,
				ScoreUploadTime:    score_uploader_time,
				ScoreReptileStatus: score_reptile_status,
			}
			scoreListTemps = append(scoreListTemps, scoreListTemp)
		}
	}
	return scoreListTemps, nil
}

func CountScoreListTemp() int {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return 0
	}
	defer db.Close()
	rows := db.QueryRow("select count(*) from score_list_temp_tbl")
	var count int
	err = rows.Scan(&count)
	if err != nil {
		log.Println("数据库查询失败 : %v", err)
	}
	return count
}

func UpdateScoreListTempStatus(href string) bool {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return false
	}
	defer db.Close()
	res, err := db.Exec("update score_list_temp_tbl set score_reptile_status = 1 where score_href = ?", href)
	if err != nil {
		log.Println("数据库更新失败 : %v", err)
		return false
	}
	rows, errs := res.RowsAffected()
	if rows < 1 {
		log.Println("更新数据库条数小于1原因是 :%v", errs)
		return false
	}
	return true
}
