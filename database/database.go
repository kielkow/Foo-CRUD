package database

import (
	"database/sql"
	"log"
	"time"
)

// DbConn to conect on database
var DbConn *sql.DB

// SetupDatabase to conect on database
func SetupDatabase() {
	var err error

	DbConn, err = sql.Open("mysql", "root:fioriferk3@tcp(127.0.0.1:3306)/inventorydb")

	if err != nil {
		log.Fatal(err)
	}

	DbConn.SetMaxOpenConns(4)
	DbConn.SetMaxIdleConns(4)
	DbConn.SetConnMaxLifetime(60 * time.Second)
}
