package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"gopkg.in/redis.v4"
	"net/http"
	"time"
)

type User struct {
	Id              int    `json:"id" orm:"column(id);auto"`
	Username        string `json:"username" orm:"column(username);size(128)"`
	Password        string `json:"password" orm:"column(password);size(128)"`
	Salt            string `json:"salt" orm:"column(salt);size(128)"`
	Email           string `json:"email" orm:"column(email);size(128)"`
	Nickname        string `json:"nickname" orm:"column(nickname);size(128)"`
	Realname        string `json:"realname" orm:"column(realname);size(128)"`
	Phone           string `json:"phone" orm:"column(phone);size(16)"`
	Idcard          string `json:"idcard" orm:"column(idcard);size(32)"`
	Register_time   string `json:"register_time" orm:"column(register_time);size(128)"`
	Register_ip     string `json:"register_ip" orm:"column(register_ip);size(128)"`
	Last_login_time string `json:"last_login_time" orm:"column(last_login_time);size(128)"`
	Last_login_ip   string `json:"last_login_ip" orm:"column(last_login_ip);size(128)"`
	Sex             int8   `json:"sex" orm:"column(sex);size(1)"`
	Birthday        string `json:"birthday" orm:"column(birthday);size(32)"`
	Province        string `json:"province" orm:"column(province);size(64)"`
	City            string `json:"city" orm:"column(city);size(64)"`
	Level           int    `json:"level" orm:"column(level);size(11)"`
	Login_count     int    `json:"login_count" orm:"column(login_count);size(11)"`
	Del             int    `json:"del" orm:"column(del);size(3)"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
	UserID   int    `json:"userID"`
	Token    string `json:"token"`
}

type CreateRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Nickname   string `json:"nickname"`
	Realname   string `json:"realname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Idcard     string `json:"idcard"`
	Sex        string `json:"sex"`
	Birthday   string `json:"birthday"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Level      string `json:"level"`
	RegisterIp string `json:"register_ip"`
}

type CreateResponse struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
}

func Login(lr *LoginRequest) (*LoginResponse, int, error) {
	// 操作Redis
	client := CreateClient()
	ListOperation(client)
	client.Close()
	username := lr.Username
	passwd := lr.Password

	if len(username) == 0 || len(passwd) == 0 {
		return nil, http.StatusBadRequest, errors.New("error:用户名或密码为空")
	}

	o := orm.NewOrm()

	user := &User{Username: username}

	err := o.Read(user, "username")
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("error:该用户不存在")
	}
	hash, _ := GeneratePassHash(passwd, user.Salt)
	if hash != user.Password {
		return nil, http.StatusBadRequest, errors.New("error:密码错误")
	}

	tokenString, err := GenerateToken(lr, user.Id, 0)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return &LoginResponse{
		Username: user.Username,
		UserID:   user.Id,
		Token:    tokenString,
	}, http.StatusOK, nil
}

// DoCreateUser: create a user
func DoCreateUser(cr *CreateRequest) (*CreateResponse, int, error) {
	// connect db
	o := orm.NewOrm()
	// check username if exist
	userNameCheck := User{Username: cr.Username}
	if len(cr.Username) == 0 || len(cr.Password) == 0 {
		return nil, http.StatusBadRequest, errors.New("用户名或密码为空")
	}
	err := o.Read(&userNameCheck, "username")
	if err == nil {
		return nil, http.StatusBadRequest, errors.New("username has already existed")
	}
	// generate salt
	saltKey, err := GenerateSalt()
	if err != nil {
		logs.Info(err.Error())

		return nil, http.StatusBadRequest, err
	}

	// generate password hash
	hash, err := GeneratePassHash(cr.Password, saltKey)
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	// create user
	user := User{}
	user.Username = cr.Username
	user.Nickname = cr.Nickname
	user.Register_time = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 03:04:05")
	user.Last_login_time = user.Register_time
	user.Password = hash
	user.Salt = saltKey
	user.Register_ip = cr.RegisterIp

	_, err = o.Insert(&user)
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	//加入es
	esClient := GetEsClient()
	put, err := esClient.Index().Index("jwt_test").Type("user").BodyJson(user).Do(context.Background())
	if err != nil {
		logs.Info("插入es失败:", err.Error)
		logs.Info("put:", put)
	}
	return &CreateResponse{
		UserID:   user.Id,
		Username: user.Username,
	}, http.StatusOK, nil

}

// redis 列表操作
func ListOperation(client *redis.Client) {
	client.RPush("fruits", "apple")
	client.LPush("fruits", "banana")
	lenght, err := client.LLen("fruits").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("lenght:", lenght)

	value, err := client.RPop("fruits").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("fruits:", value)

	value, err = client.LPop("fruits").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("fruits:", value)

}

// es 操作

func EsOperation() {

}
