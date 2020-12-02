package models

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"github.com/olivere/elastic/v6"
	"golang.org/x/crypto/scrypt"
	"gopkg.in/redis.v4"
	"io"
	"time"
)

const (
	SecretKEY              string = "JWT-Secret-Key"
	DEFAULT_EXPIRE_SECONDS int    = 600 // default expired 10 minutes
	PasswordHashBytes             = 16
)

type MyCustomClaims struct {
	UserID int `json:"userID"`
	jwt.StandardClaims
}

type JwtPayload struct {
	Username  string `json:"username"`
	UserID    int    `json:"userID"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func GenerateToken(loginInfo *LoginRequest, userId int, expireSeconds int) (tokenString string, err error) {
	if expireSeconds == 0 {
		expireSeconds = DEFAULT_EXPIRE_SECONDS
	}

	mySigningKey := []byte(SecretKEY)
	ExpiresAt := time.Now().Add(time.Second * time.Duration(expireSeconds)).Unix()

	logs.Info("token will be expired at ", time.Unix(ExpiresAt, 0))

	user := *loginInfo
	claims := MyCustomClaims{
		userId,
		jwt.StandardClaims{
			Issuer:    user.Username,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: ExpiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", errors.New("error:failed to generate token")
	}

	return tokenStr, nil
}

func GeneratePassHash(password string, salt string) (hash string, err error) {
	h, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, PasswordHashBytes)
	if err != nil {
		return "", errors.New("error: failed to generate password hash")
	}
	return fmt.Sprintf("%x", h), nil
}

// generate salt
func GenerateSalt() (salt string, err error) {
	buf := make([]byte, PasswordHashBytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", errors.New("error: failed to generate user's salt")
	}
	return fmt.Sprintf("%x", buf), nil
}

//发送邮件 TODO
func SendEmail(email string) {

}

// 获取Redis客户端

func CreateClient() *redis.Client {
	addr := beego.AppConfig.String("redis::host") + ":" + beego.AppConfig.String("redis::port")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		logs.Info("redis连接失败:", err.Error)
	}

	return client
}

// 获取es客户端
func GetEsClient() *elastic.Client {
	url := "http://" + beego.AppConfig.String("elasticsearch::host") + ":" + beego.AppConfig.String("elasticsearch::port")
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(url))
	if err != nil {
		logs.Info("elasticsearch连接失败:", err.Error)
		panic(err)
	}
	//检测es是否连接成功
	_, _, err = client.Ping(url).Do(context.Background())
	return client
}
