package foo

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

// ReportFilter struct
type ReportFilter struct {
	Name    string `json: "name"`
	Surname string `json: "surname"`
}

func handleFooReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var fooFilter ReportFilter

		err := json.NewDecoder(r.Body).Decode(&fooFilter)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		foos, err := searchFooData(fooFilter)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t := template.New("report.gotmpl")
		t, err = t.ParseFiles(path.Join("templates", "report.gotmpl"))

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var tmpl bytes.Buffer

		if len(foos) > 0 {
			err = t.Execute(&tmpl, foos)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		rdr := bytes.NewReader(tmpl.Bytes())

		w.Header().Set("Content-Disposition", "Attachement")

		http.ServeContent(w, r, "report.html", time.Now(), rdr)

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
