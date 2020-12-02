package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.Debug = true

	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		logs.Error(err.Error)
	}

	orm.RegisterModel(new(User))

	dbuser := beego.AppConfig.String("mysql::user")
	dbpass := beego.AppConfig.String("mysql::pass")
	dbhost := beego.AppConfig.String("mysql::host")
	dbport := beego.AppConfig.String("mysql::port")
	dbname := beego.AppConfig.String("mysql::db")
	dburl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", dbuser, dbpass, dbhost, dbport, dbname)
	logs.Info("connecting to mysql url:", dburl)
	err := orm.RegisterDataBase("default", "mysql", dburl)
	if err != nil {
		logs.Error(err.Error)
		panic(err.Error)
	}

}
