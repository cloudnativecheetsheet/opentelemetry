package controllers

type createTodoRequest struct {
	Content string `json:"content"`
	User_Id string `json:"user_id"`
}

type getTodosByUserRequest struct {
	User_Id string `json:"user_id"`
}

type getTodoRequest struct {
	Todo_Id string `json:"todo_id"`
}

type updateTodoRequest struct {
	Content string `json:"content"`
	User_Id string `json:"user_id"`
	Todo_Id string `json:"todo_id"`
}

type deleteTodoRequest struct {
	Todo_Id string `json:"todo_id"`
}
