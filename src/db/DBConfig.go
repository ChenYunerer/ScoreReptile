package db

import (
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

const DataSourceName = ""

var Engine *xorm.Engine

func init() {
	_engine, err := xorm.NewEngine("mysql", DataSourceName)
	if err != nil {
		panic(err)
	}
	tbMapper := names.NewSuffixMapper(names.SnakeMapper{}, "_tbl")
	_engine.SetTableMapper(tbMapper)
	Engine = _engine
}
