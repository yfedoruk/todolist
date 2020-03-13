package main

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
)

func tables(db *sql.DB) {
	files, err := filepath.Glob(basePath() + "/sql/*.sql")
	check(err)

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		check(err)

		stmt, err := db.Prepare(string(data))
		check(err)

		_, err = stmt.Exec()
		check(err)
	}
}
