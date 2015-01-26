package user

import (
	"github.com/wmbest2/rats/db"
	"log"
)

const (
	createTable = `
	CREATE TABLE users (
		id         SERIAL PRIMARY KEY,
		username   VARCHAR UNIQUE,
		password   VARCHAR,
		created_on TIMESTAMP NOT NULL DEFAULT NOW(),
	)
	`

	createTokenTable = `
	CREATE TABLE user_tokens (
		token                 VARCHAR PRIMARY KEY,
		token_encrypted       VARCHAR,
		user_id               SERIAL NOT NULL,
		created_on            TIMESTAMP NOT NULL DEFAULT NOW(),
		persistent BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (user_id) REFERENCES users(id)
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
