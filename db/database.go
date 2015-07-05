package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	Conn *sql.DB
)

func ExecOrLog(cmd string) {
	_, err := Conn.Exec(cmd)
	if err != nil {
		log.Println(err)
	}
}

func init() {
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		dburl = fmt.Sprintf("user=postgres sslmode=disable host=localhost")
	}

	var err error
	Conn, err = sql.Open("postgres", dburl)
	if err != nil {
		log.Fatal(err)
	}
}
