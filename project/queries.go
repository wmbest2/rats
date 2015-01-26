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

	createTokenTable = `
	CREATE TABLE project_tokens (
		token                 VARCHAR PRIMARY KEY,
		token_encrypted       VARCHAR,
		project_id            SERIAL UNIQUE NOT NULL,
		created_on            TIMESTAMP NOT NULL DEFAULT NOW(),
		FOREIGN KEY (project_id) REFERENCES projects(id)
	)
	`

	createProject = `
	INSERT INTO projects (name) VALUES ($1) RETURNING id
	`

	createProjectToken = `
	INSERT INTO project_tokens (project_id, token, token_encrypted) VALUES ($1, $2, $3)
	`

	updateProjectToken = `
	UPDATE project_tokens SET (token, token_encrypted, created_on) = ($2, $3, now()) where project_id = $1
	`

	findProject = `
	SELECT * FROM projects WHERE name = $1
	`

	findToken = `
	SELECT * FROM project_tokens WHERE project_id = $1
	`
)

func init() {
	log.Println("HELLO WORLD")
	_, err := db.Conn.Exec(createTable)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Conn.Exec(createTokenTable)
	if err != nil {
		log.Println(err)
	}
}
