package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func top(c *gin.Context) {
	defer LoggerAndCreateSpan(c, "TOP画面取得").End()
	generateHTML(c, "hello", "top", "layout", "top", "public_navbar", "footer")
}

func getIndex(c *gin.Context) {
	defer LoggerAndCreateSpan(c, "TODO画面取得").End()

	UserId, isExist := c.Get("UserId")
	if !isExist {
		log.Println("セッションが存在していません")
	}

	//--- UserAPI getUserByEmail への Post
	email := UserId.(string)
	jsonStr1 := `{"Email":"` + email + `"}`

	defer LoggerAndCreateSpan(c, "UserAPI /getUserByEmail にポスト").End()
	rsp1, err := otelhttp.Post(
		c.Request.Context(),
		EpUserApi+"/getUserByEmail",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr1)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp1.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp1.Body)
	var responseGetUser ResponseGetUser
	err = json.Unmarshal(byteArr, &responseGetUser)
	if err != nil {
		log.Println(err)
	}

	//--- TodoAPI getTodosByUser への Post
	user_id := strconv.Itoa(responseGetUser.ID)
	jsonStr2 := `{"user_id":"` + string(user_id) + `"}`

	defer LoggerAndCreateSpan(c, "TodoAPI /getTodosByEmail にポスト").End()
	rsp2, err := otelhttp.Post(
		c.Request.Context(),
		EpTodoAPI+"/getTodosByUser",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr2)),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp2.Body.Close()

	byteArr, _ = ioutil.ReadAll(rsp2.Body)
	var getTodosByUserresponse getTodosByUserResponse
	err = json.Unmarshal(byteArr, &getTodosByUserresponse)
	if err != nil {
		log.Println(err)
	}

	var user User
	user.Name = responseGetUser.Name
	user.Todos = getTodosByUserresponse.Todos

	defer LoggerAndCreateSpan(c, "TODO画面取得").End()
	generateHTML(c, user, "index", "layout", "private_navbar", "index", "footer")
}

func getTodoNew(c *gin.Context) {
	defer LoggerAndCreateSpan(c, "TODO作成画面取得").End()
	generateHTML(c, nil, "todoNew", "layout", "private_navbar", "todo_new", "footer")
}

func postTodoSave(c *gin.Context) {
	defer LoggerAndCreateSpan(c, "TODO保存").End()

	UserId, isExist := c.Get("UserId")
	if !isExist {
		log.Println("セッションが存在していません")
	}

	//--- UserAPI getUserByEmail への Post
	email := UserId.(string)
	jsonStr1 := `{"Email":"` + email + `"}`

	defer LoggerAndCreateSpan(c, "UserAPI /getUserByEmail にポスト").End()
	rsp1, err := otelhttp.Post(
		c.Request.Context(),
		EpUserApi+"/getUserByEmail",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr1)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp1.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp1.Body)
	var responseGetUser ResponseGetUser
	err = json.Unmarshal(byteArr, &responseGetUser)
	if err != nil {
		log.Println(err)
	}

	//--- TodoAPI createTodo への Post
	user_id := strconv.Itoa(responseGetUser.ID)
	content := c.Request.PostFormValue("content")

	defer LoggerAndCreateSpan(c, "TodoAPI /createTodo にポスト").End()
	jsonStr2 := `{"Content":"` + content + `",
	"User_Id":"` + user_id + `"}`

	rsp2, err := otelhttp.Post(
		c.Request.Context(),
		EpTodoAPI+"/createTodo",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr2)),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp2.Body.Close()

	byteArr, _ = ioutil.ReadAll(rsp2.Body)
	var getTodosByUserresponse getTodosByUserResponse
	err = json.Unmarshal(byteArr, &getTodosByUserresponse)
	if err != nil {
		log.Println(err)
	}

	defer LoggerAndCreateSpan(c, "TODO画面にリダイレクト").End()
	c.Redirect(http.StatusFound, "/menu/todos")
}

func getTodoEdit(c *gin.Context, id int) {
	defer LoggerAndCreateSpan(c, "TODO編集画面取得").End()

	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	UserId, _ := c.Get("UserId")
	//--- UserAPI getUserByEmail への Post
	email := UserId.(string)
	jsonStr1 := `{"Email":"` + email + `"}`

	defer LoggerAndCreateSpan(c, "UserAPI /getUserByEmail にポスト").End()
	rsp1, err := otelhttp.Post(
		c.Request.Context(),
		EpUserApi+"/getUserByEmail",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr1)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp1.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp1.Body)
	var responseGetUser ResponseGetUser
	err = json.Unmarshal(byteArr, &responseGetUser)
	if err != nil {
		log.Println(err)
	}

	//--- TodoAPI getTodo への Post
	todo_id := strconv.Itoa(id)
	jsonStr2 := `{"todo_id":"` + todo_id + `"}`

	defer LoggerAndCreateSpan(c, "TodoAPI /getTodo にポスト").End()
	rsp2, err := otelhttp.Post(
		c.Request.Context(),
		EpTodoAPI+"/getTodo",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr2)),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp2.Body.Close()

	byteArr, _ = ioutil.ReadAll(rsp2.Body)
	var getTodoresponse getTodoResponse
	err = json.Unmarshal(byteArr, &getTodoresponse)
	if err != nil {
		log.Println(err)
	}

	defer LoggerAndCreateSpan(c, "TODO編集画面取得").End()
	generateHTML(c, getTodoresponse, "todoEdit", "layout", "private_navbar", "todo_edit", "footer")
}

func postTodoUpdate(c *gin.Context, id int) {
	defer LoggerAndCreateSpan(c, "TODO更新").End()

	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	UserId, _ := c.Get("UserId")
	//--- UserAPI getUserByEmail への Post
	email := UserId.(string)
	jsonStr1 := `{"Email":"` + email + `"}`

	defer LoggerAndCreateSpan(c, "UserAPI /getUserByEmail にポスト").End()
	rsp1, err := otelhttp.Post(
		c.Request.Context(),
		EpUserApi+"/getUserByEmail",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr1)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp1.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp1.Body)
	var responseGetUser ResponseGetUser
	err = json.Unmarshal(byteArr, &responseGetUser)
	if err != nil {
		log.Println(err)
	}

	//--- TodoAPI updateTodo への Post
	content := c.Request.PostFormValue("content")
	user_id := strconv.Itoa(responseGetUser.ID)
	todo_id := strconv.Itoa(id)
	jsonStr2 := `{"Content":"` + content + `",
	"User_Id":"` + user_id + `",
	"Todo_Id":"` + todo_id + `"}`

	defer LoggerAndCreateSpan(c, "TodoAPI /updateTodo にポスト").End()
	rsp2, err := otelhttp.Post(
		c.Request.Context(),
		EpTodoAPI+"/updateTodo",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr2)),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp2.Body.Close()

	byteArr, _ = ioutil.ReadAll(rsp2.Body)
	var updateTodoresponse updateTodoResponse
	err = json.Unmarshal(byteArr, &updateTodoresponse)
	if err != nil {
		log.Println(err)
	}

	defer LoggerAndCreateSpan(c, "TODO画面にリダイレクト").End()
	c.Redirect(http.StatusFound, "/menu/todos")
}

func getTodoDelete(c *gin.Context, id int) {
	defer LoggerAndCreateSpan(c, "TODO削除").End()

	//--- TodoAPI deleteTodo への Post
	todo_id := strconv.Itoa(id)
	jsonStr1 := `{"todo_id":"` + todo_id + `"}`

	defer LoggerAndCreateSpan(c, "TodoAPI /deleteTodo にポスト").End()
	rsp1, err := otelhttp.Post(
		c.Request.Context(),
		EpTodoAPI+"/deleteTodo",
		"application/json",
		bytes.NewBuffer([]byte(jsonStr1)),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp1.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp1.Body)
	var deleteTodoresponse deleteTodoResponse
	err = json.Unmarshal(byteArr, &deleteTodoresponse)
	if err != nil {
		log.Println(err)
	}

	defer LoggerAndCreateSpan(c, "TODO画面にリダイレクト").End()
	c.Redirect(http.StatusFound, "/menu/todos")
}
