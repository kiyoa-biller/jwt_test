package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"jwt_demo/models"

	"net/http"
)

type UserController struct {
	beego.Controller
}

func (u *UserController) unmarshalPayload(v interface{}) error {
	body := u.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, v); err != nil {
		logs.Error("unmarshal payload of %s error: %s", u.Ctx.Request.URL.Path, err)
	}
	return nil
}

func (u *UserController) respond(code int, message string, data ...interface{}) {
	u.Ctx.Output.SetStatus(code)
	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}
	u.Data["json"] = struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}{
		Code:    code,
		Message: message,
		Data:    d,
	}
	u.ServeJSON()
}

func (u *UserController) Login() {

	lr := new(models.LoginRequest)

	if err := u.unmarshalPayload(lr); err != nil {
		u.respond(http.StatusBadRequest, err.Error())
	}
	lrs, statusCode, err := models.Login(lr)
	if err != nil {
		u.respond(statusCode, err.Error())
		return
	}
	u.Ctx.Output.Header("Authorization", lrs.Token) // set token into header
	u.respond(http.StatusOK, "", lrs)

}

// CreateUser creates a user
func (u *UserController) CreateUser() {
	cu := new(models.CreateRequest)

	if err := u.unmarshalPayload(cu); err != nil {
		u.respond(http.StatusBadRequest, err.Error())
	}
	//获取IP地址
	cu.RegisterIp = u.Ctx.Request.RemoteAddr
	createUser, statusCode, err := models.DoCreateUser(cu)
	if err != nil {
		u.respond(statusCode, err.Error())
		return
	}

	u.respond(http.StatusOK, "", createUser)
}
