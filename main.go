package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pluralsight/inventoryservice/database"
	"github.com/pluralsight/inventoryservice/foo"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	foo.SetupRoutes(apiBasePath)
	http.ListenAndServe(":3333", nil)
}
