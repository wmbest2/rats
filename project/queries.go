package project

import (
	"github.com/wmbest2/rats/db"
	"log"
)

const (
	createTable = `
	CREATE TABLE projects (
		id         SERIAL PRIMARY KEY,
		name       VARCHAR UNIQUE,
		created_on TIMESTAMP NOT NULL DEFAULT NOW()
	)
	`

	createProject = `
	INSERT INTO projects (name) VALUES ($1) RETURNING id
	`

	findProject = `
	SELECT * FROM projects WHERE name = $1
	`
)

func init() {
	_, err := db.Conn.Exec(createTable)
	if err != nil {
		log.Println(err)
	}
}
