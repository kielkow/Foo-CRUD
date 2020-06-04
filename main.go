package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pluralsight/inventoryservice/database"
	"github.com/pluralsight/inventoryservice/foo"
	"github.com/pluralsight/inventoryservice/receipt"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()

	receipt.SetupRoutes(apiBasePath)
	foo.SetupRoutes(apiBasePath)

	log.Fatal(http.ListenAndServe(":3333", nil))
}
