package project

import (
	"github.com/wmbest2/rats/db"
	"log"
)

const (
	createTable = `
	CREATE TABLE projects (
		id         SERIAL PRIMARY KEY,
		name       VARCHAR,
		created_on TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT uni_name UNIQUE (name)
	)
	`

	createTokenTable = `
	CREATE TABLE project_tokens (
		token                 VARCHAR PRIMARY KEY,
		token_encrypted       VARCHAR,
		project_id            SERIAL UNIQUE NOT NULL,
		created_on            TIMESTAMP NOT NULL DEFAULT NOW(),
		FOREIGN KEY (project_id) REFERENCES project(id)
	)
	`
)

func init() {
	_, err := db.Conn.Exec(createTable)
	if err != nil {
		log.Println(err)
	}

	_, err := db.Conn.Exec(createTokenTable)
	if err != nil {
		log.Println(err)
	}
}
