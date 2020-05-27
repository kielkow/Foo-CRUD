package main

import (
	"net/http"

	"github.com/pluralsight/inventoryservice/foo"
)

const apiBasePath = "/api"

func main() {
	foo.SetupRoutes(apiBasePath)
	http.ListenAndServe(":3333", nil)
}
