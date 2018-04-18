package main

import (
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"confirmOffert", "POST", "/offert/{id}", confirmOffert})
}

const (
	FmDatabase = "G-Smart"
	OrderTable = "www_order"
	LogField   = "Logg"
)

func confirmOffert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fms := filemaker.NewServer(settings.FileMakerHost, settings.FileMakerUser, settings.FileMakerPassword)

	params := make([]filemaker.SearchParam, 0)
	params = append(params, filemaker.SearchParam{Op: filemaker.Equal, Name: "wwwId", Value: vars["id"]})
	reqs, n, err := fms.Get(FmDatabase, OrderTable, params, 0, 10)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n != 1 {
		// We got more than one, something is wrong
		logger.Error(vars["id"], " found ", n)
		http.Error(w, "Unable to find record", http.StatusInternalServerError)
		return
	}
	id := reqs[0]["recid"].String()
	fmData := make(map[string]string)
	fmData[LogField] = reqs[0][LogField].String() + "Server: Offert ACCEPTERAD " + time.Now().Format("2006-01-02 15.04.05") + "\n"
	_, err = fms.EditRow(FmDatabase, OrderTable, id, fmData)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		OutputText string
	}{OutputText: "Offerten är nu accepterad. Vi tackar för ditt förtroende."}
	dir, _ := getDir()
	t, _ := template.ParseFiles(dir + "html/confirmOffert.html")
	if err := t.Execute(w, data); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
