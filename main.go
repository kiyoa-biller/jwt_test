package main

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	_ "jwt_demo/routers"
)

// //测试MySQL连接
// func init() {
// 	orm.Debug = true
// 	if err := orm.RegisterDriver("mysql"); err != nil {
// 		fmt.Println("数据库连接失败")
// 	}
// }
func main() {
	beego.Run()
}
