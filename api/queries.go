package api

import (
	"github.com/wmbest2/rats/db"
	"log"
)

const (
	createTokenTable = `
	CREATE TABLE api_tokens (
		id                    SERIAL PRIMARY KEY,
		token                 VARCHAR ,
		token_encrypted       VARCHAR,
		type                  INTEGER,
		parent_id             SERIAL UNIQUE,
		created_on            TIMESTAMP NOT NULL DEFAULT NOW()
	)
	`

	createToken = `
	INSERT INTO api_tokens (type, parent_id, token, token_encrypted) VALUES ($1, $2, $3, $4)
	`

	updateToken = `
	UPDATE api_tokens SET (token, token_encrypted, created_on) = ($2, $3, now()) where parent_id = $1
	`

	findToken = `
	SELECT * FROM api_tokens WHERE parent_id = $1
	`

	findEncryptedToken = `
	SELECT id FROM api_tokens WHERE token_encrypted = $1
	`
)

func init() {
	_, err := db.Conn.Exec(createTokenTable)
	if err != nil {
		log.Println(err)
	}
}
