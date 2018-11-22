package db

import (
	"ScoreReptile/src/model"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func InsertScorePictureInfo(socrePictureInfo model.ScorePictureInfo) error {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return err
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO score_picture_info_tbl(score_id, score_name, score_href, score_picture_index, score_picture_href) VALUES (?, ?, ?, ?, ?)",
		socrePictureInfo.ScoreId, socrePictureInfo.ScoreName, socrePictureInfo.ScoreHref, socrePictureInfo.ScorePictureIndex, socrePictureInfo.ScorePictureHref)
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

func CountScorePictureInfo(scoreHref string) int64 {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return 0
	}
	defer db.Close()
	rows := db.QueryRow("select count(*) from score_picture_info_tbl where score_href = ?", scoreHref)
	var count int64
	err = rows.Scan(&count)
	if err != nil {
		log.Println("数据库查询失败 : %v", err)
		return 0
	}
	return count
}
