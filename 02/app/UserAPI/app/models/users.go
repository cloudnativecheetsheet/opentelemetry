package models

import (
	"log"
	"time"
	"userapi/app/utils"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	PassWord  string
	CreatedAt time.Time
}

func (u *User) CreateUser(c *gin.Context) (err error) {
	defer utils.LoggerAndCreateSpan(c, "CRUD : CreateUser")

	cmd := `insert into users (
		uuid,
		name,
		email,
		password,
		created_at) values ($1, $2, $3, $4, $5)`

	_, err = Db.Exec(cmd,
		createUUID(c),
		u.Name,
		u.Email,
		Encrypt(c, u.PassWord),
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUser(c *gin.Context, id int) (user User, err error) {
	defer utils.LoggerAndCreateSpan(c, "CRUD : GetUser")

	user = User{}
	cmd := `select id, uuid, name, email, password, created_at
	from users where id = $1`
	err = Db.QueryRow(cmd, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt)

	return user, err
}

func (u *User) UpdateUser(c *gin.Context) (err error) {
	defer utils.LoggerAndCreateSpan(c, "CRUD : UpdateUser").End()

	cmd := `update users set name = $1, email = $2 where id = $3`
	_, err = Db.Exec(cmd, u.Name, u.Email, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (u *User) DeleteUser(c *gin.Context) (err error) {
	defer utils.LoggerAndCreateSpan(c, "CRUD : DeleteUser").End()

	cmd := `delete from users where id = $1`
	_, err = Db.Exec(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUserByEmail(c *gin.Context, email string) (user User, err error) {
	defer utils.LoggerAndCreateSpan(c, "CRUD : GetUserByEmail").End()

	user = User{}
	cmd := `select id, uuid, name, email, password, created_at
	from users where email = $1`
	err = Db.QueryRow(cmd, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt)

	// for CNDT2022
	// time.Sleep(100 * time.Millisecond)

	return user, err
}
