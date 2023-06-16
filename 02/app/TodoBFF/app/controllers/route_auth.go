package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func getSignup(c *gin.Context) {
	defer LoggerAndCreateSpan(c, "ユーザ登録画面取得").End()
	generateHTML(c, nil, "signup", "layout", "signup", "public_navbar", "footer")
}

func postSignup(c *gin.Context) {
	defer LoggerAndCreateSpan(c, "ユーザ登録").End()
	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	//--- UserAPI createUser への Post
	name := c.Request.PostFormValue("name")
	email := c.Request.PostFormValue("email")
	password := c.Request.PostFormValue("password")

	jsonStr := `{"Name":"` + name + `",
	"Email":"` + email + `",
	"PassWord":"` + password + `"}`

	defer LoggerAndCreateSpan(c, "UserAPI /createUser にポスト").End()
	rsp, err := otelhttp.Post(
		c.Request.Context(),
		EpUserApi+"/createUser",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)
	log.Println(string(body))

	UserId := email
	login(c, UserId)

	defer LoggerAndCreateSpan(c, "TODO画面にリダイレクト").End()
	c.Redirect(http.StatusMovedPermanently, "/menu/todos")
}

func getLogin(c *gin.Context) {
	defer LoggerAndCreateSpan(c, "ログイン画面取得").End()
	generateHTML(c, nil, "login", "layout", "login", "public_navbar", "footer")
}

func postLogin(c *gin.Context) {
	defer LoggerAndCreateSpan(c, "ログイン").End()
	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	//--- UserAPI getUserByEmail への Post
	email := c.Request.PostFormValue("email")
	jsonStr := `{"Email":"` + email + `"}`

	LoggerAndCreateSpan(c, "UserAPI /getUserByEmail にポスト").End()
	rsp, err := otelhttp.Post(
		c.Request.Context(),
		EpUserApi+"/getUserByEmail",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp.Body)
	var responseGetUser ResponseGetUser
	err = json.Unmarshal(byteArr, &responseGetUser)
	if err != nil {
		log.Println(err)
	}

	//--- UserAPI encrypt への Post
	password := c.Request.PostFormValue("password")
	jsonStr = `{"PassWord":"` + password + `"}`

	defer LoggerAndCreateSpan(c, "UserAPI /encrypt にポスト").End()
	rsp, err = otelhttp.Post(
		c.Request.Context(),
		EpUserApi+"/encrypt",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()

	byteArr, _ = ioutil.ReadAll(rsp.Body)
	var responseEncrypt ResponseEncrypt
	err = json.Unmarshal(byteArr, &responseEncrypt)
	if err != nil {
		log.Println(err)
	}

	if responseGetUser.ID == 0 {
		log.Println("ユーザがいません")

		defer LoggerAndCreateSpan(c, "ログイン画面にリダイレクト").End()
		c.Redirect(http.StatusFound, "/login")
	} else if responseEncrypt.PassWord == responseGetUser.PassWord {
		UserId := c.PostForm("email")
		login(c, UserId)

		defer LoggerAndCreateSpan(c, "TODO画面にリダイレクト").End()
		c.Redirect(http.StatusMovedPermanently, "/menu/todos")
	} else {
		log.Println("PW が間違っています")

		defer LoggerAndCreateSpan(c, "ログイン画面にリダイレクト").End()
		c.Redirect(http.StatusFound, "/login")
	}
}

func getLogout(c *gin.Context) {
	defer LoggerAndCreateSpan(c, "ログアウト").End()
	logout(c)

	defer LoggerAndCreateSpan(c, "TOP画面にリダイレクト").End()
	c.Redirect(http.StatusMovedPermanently, "/")
}

func login(c *gin.Context, UserId string) {
	defer LoggerAndCreateSpan(c, "ログイン処理...").End()

	session := sessions.Default(c)

	defer LoggerAndCreateSpan(c, "セッション設定").End()
	session.Set("UserId", UserId)

	defer LoggerAndCreateSpan(c, "セッション保存").End()
	session.Save()

	defer LoggerAndCreateSpan(c, "ログイン完了").End()
}

func logout(c *gin.Context) {
	defer LoggerAndCreateSpan(c, "ログアウト処理...").End()

	session := sessions.Default(c)

	defer LoggerAndCreateSpan(c, "セッションクリア").End()

	session.Clear()

	defer LoggerAndCreateSpan(c, "セッション保存").End()
	session.Save()

	defer LoggerAndCreateSpan(c, "ログアウト完了").End()
}
