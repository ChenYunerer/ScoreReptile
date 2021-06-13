package db

import (
	"ScoreReptile/src/model"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func InsertScoreBaseInfo(scoreBaseInfo model.ScoreBaseInfo) error {
	db, err := sql.Open("mysql", DataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return err
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO score_base_info_tbl(score_id, score_category, score_name, score_href, score_singer, score_author, score_word_writer, score_song_writer, score_format, score_origin, score_uploader, score_uploader_time, score_view_count, score_cover_picture, score_picture_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		scoreBaseInfo.ScoreId, scoreBaseInfo.ScoreCategory, scoreBaseInfo.ScoreName, scoreBaseInfo.ScoreHref, scoreBaseInfo.ScoreSinger, scoreBaseInfo.ScoreAuthor, scoreBaseInfo.ScoreWordWriter, scoreBaseInfo.ScoreSongWriter, scoreBaseInfo.ScoreFormat, scoreBaseInfo.ScoreOrigin, scoreBaseInfo.ScoreUploader, scoreBaseInfo.ScoreUploadTime, scoreBaseInfo.ScoreViewCount, scoreBaseInfo.ScoreCoverPicture, scoreBaseInfo.ScorePictureCount)
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

func IsScoreBaseInfoExist(href string) bool {
	db, err := sql.Open("mysql", DataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return false
	}
	defer db.Close()
	rows := db.QueryRow("select count(*) from score_base_info_tbl where score_href = ?", href)
	var count int64
	err = rows.Scan(&count)
	if err != nil {
		log.Println("数据库查询失败 : %v", err)
	}
	return count > 0
}

func UpdateScoreBaseInfoId(href string, id int) bool {
	db, err := sql.Open("mysql", DataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return false
	}
	defer db.Close()
	res, err := db.Exec("update score_base_info_tbl set score_id = ? where score_href = ?", id, href)
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

func GetScoreBaseInfo(count int) ([]model.ScoreBaseInfo, error) {
	scoreBaseInfos := make([]model.ScoreBaseInfo, 0)
	db, err := sql.Open("mysql", DataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select score_id, score_name, score_href from score_base_info_tbl where score_picture_count = 0 limit ?", count)
	if err != nil {
		log.Println("数据查询库失败 : %v", err)
		return nil, err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		i++
		log.Println("处理数据库数据", i)
		var score_id int
		var score_name, score_href string
		err = rows.Scan(&score_id, &score_name, &score_href)
		if err != nil {
			log.Println("数据查询对象映射失败 : %v", err)
		} else {
			scoreBaseInfo := model.ScoreBaseInfo{
				ScoreId:   score_id,
				ScoreName: score_name,
				ScoreHref: score_href,
			}
			scoreBaseInfos = append(scoreBaseInfos, scoreBaseInfo)
		}
	}
	return scoreBaseInfos, nil
}

func GetUnCountPicScoreBaseInfo() ([]model.ScoreBaseInfo, error) {
	scoreBaseInfos := make([]model.ScoreBaseInfo, 0)
	db, err := sql.Open("mysql", DataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select score_id, score_name, score_href from score_base_info_tbl where score_picture_count = 0")
	if err != nil {
		log.Println("数据查询库失败 : %v", err)
		return nil, err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		i++
		log.Println("处理数据库数据", i)
		var score_id int
		var score_name, score_href string
		err = rows.Scan(&score_id, &score_name, &score_href)
		if err != nil {
			log.Println("数据查询对象映射失败 : %v", err)
		} else {
			scoreBaseInfo := model.ScoreBaseInfo{
				ScoreId:   score_id,
				ScoreName: score_name,
				ScoreHref: score_href,
			}
			scoreBaseInfos = append(scoreBaseInfos, scoreBaseInfo)
		}
	}
	return scoreBaseInfos, nil
}

func UpdateScoreBaseInfoPictureCount(href string, count int64) bool {
	db, err := sql.Open("mysql", DataSourceName)
	if err != nil {
		log.Println("打开数据库失败 : %v", err)
		return false
	}
	defer db.Close()
	res, err := db.Exec("update score_base_info_tbl set score_picture_count = ? where score_href = ?", count, href)
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
