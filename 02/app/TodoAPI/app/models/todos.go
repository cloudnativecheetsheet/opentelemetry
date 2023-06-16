package models

import (
	"log"
	"time"
	"todoapi/app/utils"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID        int
	Content   string
	UserID    int
	CreatedAt time.Time
}

func CreateTodo(c *gin.Context, content string, user_id string) (err error) {
	defer utils.LoggerAndCreateSpan(c, "CRUD : CreateTodo").End()

	cmd := `insert into todos (
				content,
				user_id,
				created_at) values ($1, $2, $3)`

	_, err = Db.Exec(cmd, content, user_id, time.Now())

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func GetTodo(c *gin.Context, todo_id string) (todo Todo, err error) {
	defer utils.LoggerAndCreateSpan(c, "CRUD : GetTodo").End()

	cmd := `select id, content, user_id, created_at from todos
	where id = $1`
	todo = Todo{}

	err = Db.QueryRow(cmd, todo_id).Scan(
		&todo.ID,
		&todo.Content,
		&todo.UserID,
		&todo.CreatedAt)

	return todo, err
}

func GetTodos(c *gin.Context) (todos []Todo, err error) {
	defer utils.LoggerAndCreateSpan(c, "CRUD : GetTodos").End()

	cmd := `select id, content, user_id, created_at from todos`
	rows, err := Db.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		todos = append(todos, todo)
	}
	rows.Close()

	return todos, err
}

func GetTodosByUser(c *gin.Context, user_id string) (todos []Todo, err error) {
	defer utils.LoggerAndCreateSpan(c, "CRUD : GetTodosByUser").End()

	cmd := `select id, content, user_id, created_at from todos
	where user_id = $1`

	rows, err := Db.Query(cmd, user_id)
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		todos = append(todos, todo)
	}
	rows.Close()

	return todos, err
}

func UpdateTodo(c *gin.Context, content string, user_id string, todo_id string) error {
	defer utils.LoggerAndCreateSpan(c, "CRUD : UpdateTodo").End()

	cmd := `update todos set content = $1, user_id = $2 
	where id = $3`
	_, err = Db.Exec(cmd, content, user_id, todo_id)

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func DeleteTodo(c *gin.Context, todo_id string) error {
	defer utils.LoggerAndCreateSpan(c, "CRUD : DeleteTodo").End()

	cmd := `delete from todos where id = $1`
	_, err = Db.Exec(cmd, todo_id)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}
