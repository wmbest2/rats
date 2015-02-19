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
		is_admin   BOOLEAN DEFAULT FALSE,
		created_on TIMESTAMP NOT NULL DEFAULT NOW()
	)
	`

	createOauthTable = `
	CREATE TABLE oauth_ids (
		oauth_id              VARCHAR PRIMARY KEY,
		service               VARCHAR,
		user_id               SERIAL NOT NULL,
		created_on            TIMESTAMP NOT NULL DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
	`

	createUser = `
	INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3) RETURNING id
	`

	updateUserPassword = `
	UPDATE users SET (password) = ($2) where id = $1
	`

	updateUserRole = `
	UPDATE users SET (is_admin) = ($2) where id = $1
	`

	getAllUsers = `
	SELECT * FROM users
	`

	findUser = `
	SELECT * FROM users WHERE username = $1
	`
)

func init() {
	_, err := db.Conn.Exec(createTable)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Conn.Exec(createOauthTable)
	if err != nil {
		log.Println(err)
	}

	New("admin", "admin", true)
}
