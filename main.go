package main

import (
	"database/sql"
)

var db *sql.DB

func main() {
	if db != nil {
		defer db.Close()
	}
}
