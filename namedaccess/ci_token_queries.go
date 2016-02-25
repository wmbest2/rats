package namedaccess

import (
	"github.com/wmbest2/rats/db"
	"log"
)

const (
	createAccessTable = `
	CREATE TABLE namedAccess (
		id         SERIAL PRIMARY KEY,
		name       VARCHAR,
		project_id SERIAL,
		created_on TIMESTAMP NOT NULL DEFAULT NOW()
	)
	`

	createNamedAccess = `
	INSERT INTO namedAccess (name, project_id) VALUES ($1, $2) RETURNING id
	`

	findNamedAccess = `
	SELECT * FROM namedAccess WHERE name = $1
	`

	findNamedAccessByProject = `
	SELECT * FROM namedAccess WHERE project_id = $1
	`
)

func init() {
	_, err := db.Conn.Exec(createAccessTable)
	if err != nil {
		log.Println(err)
	}
}
