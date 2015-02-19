package api

import (
	"github.com/wmbest2/rats/db"
	"log"
)

const (
	createTokenTable = `
	CREATE TABLE api_tokens (
		token                 VARCHAR PRIMARY KEY,
		token_encrypted       VARCHAR,
		type                  INTEGER,
		id                    SERIAL NOT NULL,
		created_on            TIMESTAMP NOT NULL DEFAULT NOW()
	)
	`

	createToken = `
	INSERT INTO api_tokens (type, id, token, token_encrypted) VALUES ($1, $2, $3, $4)
	`

	updateToken = `
	UPDATE api_tokens SET (token, token_encrypted, created_on) = ($2, $3, now()) where id = $1
	`

	findToken = `
	SELECT * FROM api_tokens WHERE id = $1
	`
)

func init() {
	_, err := db.Conn.Exec(createTokenTable)
	if err != nil {
		log.Println(err)
	}
}
