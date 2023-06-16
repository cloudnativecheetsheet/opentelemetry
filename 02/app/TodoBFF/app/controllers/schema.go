package controllers

import "time"

type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	PassWord  string
	CreatedAt time.Time
	Todos     []Todo
}

type Todo struct {
	ID        int
	Content   string
	UserID    int
	CreatedAt time.Time
}

type Todos struct {
	Todos []Todo
}

type getTodosByUserResponse struct {
	Todos []Todo `json:"todos"`
}

type getTodoResponse struct {
	ID        int       `json:"ID"`
	Content   string    `json:"Content"`
	UserID    int       `json:"UserID"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type updateTodoResponse struct {
	Content string `json:"Content"`
}

type deleteTodoResponse struct {
	ResultCode string `json:"resultCode"`
}

type ResponseGetUser struct {
	ID        int    `json:"ID"`
	UUID      string `json:"UUID"`
	Name      string `json:"Name"`
	Email     string `json:"Email"`
	PassWord  string `json:"PassWord"`
	CreatedAt string `json:"CreatedAt"`
}

type ResponseEncrypt struct {
	PassWord string `json:"PassWord"`
}
