package model

import "time"

type NeteaseMusicTopListDetailInfo struct {
	Id                int       `json:"id"`
	SongName          string    `xorm:"song_name" json:"songName"`                     //歌曲名称
	Sort              int       `xorm:"sort" json:"sort"`                              //歌曲排名 从0开始
	Href              string    `xorm:"href" json:"href"`                              //网易云歌曲访问地址（相对地址）
	TopListId         int       `xorm:"top_list_id" json:"topListId"`                  //网易云榜单id
	TopListName       string    `xorm:"top_list_name" json:"topListName"`              //网易云榜单名称
	TopListUpdateTime time.Time `xorm:"top_list_update_time" json:"topListUpdateTime"` //网易云榜单更新时间
	CreateTime        time.Time `xorm:"create_time" json:"createTime"`                 //创建时间
	UpdateTime        time.Time `xorm:"update_time" json:"updateTime"`                 //更新时间
	TopTaskId         string    `xorm:"top_task_id" json:"topTaskId"`                  //所属顶层父任务ID
	TaskId            string    `xorm:"task_id" json:"taskId"`                         //所属任务ID
}
