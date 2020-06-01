package database

import (
	"database/sql"
	"log"
)

// DbConn to conect on database
var DbConn *sql.DB

// SetupDatabase to conect on database
func SetupDatabase() {
	var err error

	DbConn, err = sql.Open("mysql", "root:password123@tcp(127.0.0.1:3006)/inventorydb")

	if err != nil {
		log.Fatal(err)
	}
}
