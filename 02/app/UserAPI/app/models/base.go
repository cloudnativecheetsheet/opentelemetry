package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB
var err error

func init() {
	Db, err = sql.Open("postgres", "host=postgresql.todo.svc.cluster.local port=5432 user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Println(err)
	}

	cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id serial PRIMARY KEY,
		uuid text NOT NULL UNIQUE,
		name text,
		email text,
		password text,
		created_at timestamp)`, "users")

	Db.Exec(cmdU)
}

func createUUID(c *gin.Context) (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(c *gin.Context, plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
