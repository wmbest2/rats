package user

import (
	"github.com/wmbest2/rats/db"
	"log"
)

const (
	createTable = `
	CREATE TABLE users (
		id         SERIAL PRIMARY KEY,
		username   VARCHAR,
		password   VARCHAR
	)
	`

	createCITokenTable = `
	CREATE TABLE ci_token (
		token                 UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		created_on            TIMESTAMP NOT NULL DEFAULT NOW()
	)
	`

	generateCIToken = `
		INSERT INTO ci_token VALUES (); 
		CREATE RULE no_insert AS ON INSERT TO ci_token DO INSTEAD NOTHING; 
		CREATE RULE no_delete AS ON DELETE TO ci_token DO INSTEAD NOTHING; 
	`

	createTokenTable = `
	CREATE TABLE user_tokens (
		token                 VARCHAR PRIMARY KEY,
		token_encrypted       VARCHAR,
		user_id               SERIAL NOT NULL,
		created_on            TIMESTAMP NOT NULL DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
	`
)

func init() {
	_, err := db.Conn.Exec(createTable)
	if err != nil {
		log.Println(err)
	}

	_, err := db.Conn.Exec(createCITokenTable)
	if err != nil {
		log.Println(err)
	}

	_, err := db.Conn.Exec(generateCIToken)
	if err != nil {
		log.Println(err)
	}

	_, err := db.Conn.Exec(createTokenTable)
	if err != nil {
		log.Println(err)
	}
}
