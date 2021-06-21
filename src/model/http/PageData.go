package http

type PageData struct {
	PageNo    int64       `json:"pageNo"`
	PageSize  int64       `json:"pageSize"`
	Total     int64       `json:"total"`
	TotalPage int64       `json:"totalPage"`
	ListData  interface{} `json:"listData"`
}
