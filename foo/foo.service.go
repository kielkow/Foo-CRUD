package foo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/pluralsight/inventoryservice/cors"
	"golang.org/x/net/websocket"
)

const foosBasePath = "foos"

// SetupRoutes function
func SetupRoutes(apiBasePath string) {
	handleFoos := http.HandlerFunc(foosHandler)
	handleFoo := http.HandlerFunc(fooHandler)

	reportHandler := http.HandlerFunc(handleFooReport)

	http.Handle("/websocket", websocket.Handler(fooSocket))

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, foosBasePath), cors.Middleware(handleFoos))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, foosBasePath), cors.Middleware(handleFoo))

	http.Handle(fmt.Sprintf("%s/%s/reports", apiBasePath, foosBasePath), cors.Middleware(reportHandler))
}

func foosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fooList, err := getFooList()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		foosJSON, err := json.Marshal(fooList)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(foosJSON)

	case http.MethodPost:
		var newFoo Foo
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newFoo)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newFoo.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = insertFoo(newFoo)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return

	case http.MethodOptions:
		return
	}
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "foos/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	foo, err := getFoo(productID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if foo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		fooJSON, err := json.Marshal(foo)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(fooJSON)

	case http.MethodPut:
		var updatedFoo Foo

		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedFoo)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if updatedFoo.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = updateFoo(updatedFoo)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	case http.MethodDelete:
		err = removeFoo(productID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
