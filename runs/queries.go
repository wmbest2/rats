package runs

import (
	"log"

	"github.com/wmbest2/rats/db"
)

const (
	createRunTable = `
	CREATE TABLE runs (
		id         SERIAL PRIMARY KEY,
		token_id   SERIAL,
		project    SERIAL,
		commit     VARCHAR,
		message    VARCHAR,
		timestamp  TIMESTAMP NOT NULL DEFAULT NOW(),
		time       INTEGER,
		success    BOOLEAN,
		FOREIGN KEY (token_id) REFERENCES api_tokens(id)
	)
	`

	createRun = `
	INSERT INTO runs (token_id, project, timestamp) VALUES ($1, $2, $3) RETURNING id
	`

	saveRun = `
	UPDATE runs SET (commit, message, time, success) VALUES ($2, $3, $4, $5) where id = $1
	`

	findRun = `
	SELECT * FROM runs where id = $1
	`

	createArtifactTable = `
	CREATE TABLE run_artifacts (
		id         SERIAL PRIMARY KEY,
		suite_id   SERIAL,
		name       VARCHAR UNIQUE,
		data       BYTEA,
		FOREIGN KEY (suite_id) REFERENCES suites(id)
	)
	`

	createSuiteTable = `
	CREATE TABLE suites (
		id         SERIAL PRIMARY KEY,
		run_id     SERIAL,
		tests      INTEGER,
		failures   INTEGER,
		errors     INTEGER,
		skipped    INTEGER,
		hostname   VARCHAR,
		time       INTEGER,
		name       VARCHAR,
		stdout     VARCHAR,
		stderr     VARCHAR,
		FOREIGN KEY (run_id) REFERENCES runs(id)
	)
	`

	createSuitePropertiesTable = `
	CREATE TABLE suite_properties (
		id         SERIAL PRIMARY KEY,
		name       VARCHAR UNIQUE,
		value      VARCHAR,
		FOREIGN KEY (suite_id) REFERENCES suites(id)
	)
	`

	createCaseTable = `
	CREATE TABLE cases (
		id         SERIAL PRIMARY KEY,
		suite_id   SERIAL,
		classname  VARCHAR,
		name       VARCHAR,
		status     VARCHAR,
		assertions VARCHAR,
		time       INTEGER,
		failed     BOOLEAN,
		skipped    BOOLEAN,
		FOREIGN KEY (suite_id) REFERENCES suites(id)
	)
	`

	createStackTable = `
	CREATE TABLE stacktraces (
		id         SERIAL PRIMARY KEY,
		case_id   SERIAL,
		type      INTEGER,
		stack     VARCHAR,
		FOREIGN KEY (case_id) REFERENCES cases(id)
	)
	`
)

func init() {
	_, err := db.Conn.Exec(createRunTable)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Conn.Exec(createArtifactTable)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Conn.Exec(createSuiteTable)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Conn.Exec(createSuitePropertiesTable)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Conn.Exec(createCaseTable)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Conn.Exec(createStackTable)
	if err != nil {
		log.Println(err)
	}

}
