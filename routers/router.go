package routers

import (
	"github.com/astaxie/beego"
	"jwt_demo/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/user",
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
			beego.NSRouter("/register", &controllers.UserController{}, "post:CreateUser"),
		),
	)
	beego.AddNamespace(ns)
}
