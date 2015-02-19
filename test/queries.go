package test

const (
	createRunTable = `
	CREATE TABLE runs (
		id         SERIAL PRIMARY KEY,
		token_id   SERIAL,
		project    VARCHAR,
		commit     VARCHAR,
		message    VARCHAR,
		timestamp  TIMESTAMP NOT NULL DEFAULT NOW(),
		time       INTEGER,
		success    BOOLEAN
		FOREIGN KEY (token_id) REFERENCES api_tokens(id)
	)
	`

	createArtifactTable = `
	CREATE TABLE suite_properties (
		id         SERIAL PRIMARY KEY,
		run_id     SERIAL,
		name       VARCHAR UNIQUE,
		data       BLOB,
		FOREIGN KEY (run_id) REFERENCES runs(id)
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
