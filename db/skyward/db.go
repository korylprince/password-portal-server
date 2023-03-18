package skyward

import (
	"database/sql"
	"fmt"
	"time"
)

// DB represents a connection to an Skyward server
type DB struct {
	db *sql.DB
}

// New returns a new *DB with the given parameters
func New(dsn string) (*DB, error) {
	db, err := sql.Open("odbc", dsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to open database: %v", err)
	}
	return &DB{db: db}, nil
}

// GetID returns the id for the student matching the given information,
// an empty string if the information doesn't match, or an error if one occurred
func (db *DB) GetID(firstName, lastName, ssn string, birthDate time.Time) (string, error) {
	row := db.db.QueryRow(`
		SELECT
			student."OTHER-ID" AS ID
		FROM PUB.NAME AS name
		INNER JOIN PUB."STUDENT" AS student ON
				name."NAME-ID" = student."NAME-ID"

		INNER JOIN PUB."STUDENT-ENTITY" as sentity ON
				sentity."STUDENT-ID" = student."STUDENT-ID"

		INNER JOIN PUB."ENTITY" as entity ON
				entity."ENTITY-ID" = sentity."ENTITY-ID"
		WHERE
			sentity."STUDENT-STATUS" = 'A' AND
			student."GRAD-YR" >= entity."SCHOOL-YEAR" AND
			name."FIRST-NAME" = ? AND
			name."LAST-NAME" = ? AND
			TO_CHAR(name."BIRTHDATE", 'YYYY-MM-DD') = ? AND
			SUBSTR(name."FEDERAL-ID-NO", 6, 4) = ?

		WITH (NOLOCK)
	`,
		firstName,
		lastName,
		birthDate.Format("2006-01-02"),
		ssn,
	)

	var id string

	err := row.Scan(
		&id,
	)

	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", fmt.Errorf("Unable to query student: %v", err)
	}

	return id, nil
}
