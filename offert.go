package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johansundell/fmp-json/filemaker"
)

func init() {
	routes = append(routes, Route{"offert", "GET", "/offert/{id}", offertHandler})
}

func offertHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars

	fms := filemaker.NewServer(settings.FileMakerHost, settings.FileMakerUser, settings.FileMakerPassword)

	params := make([]filemaker.SearchParam, 0)
	params = append(params, filemaker.SearchParam{Op: filemaker.Equal, Name: "wwwId", Value: vars["id"]})
	reqs, n, err := fms.Get("G-Smart", "www_order", params, 0, 10)
	if err != nil {
		logger.Error(err)
		return
	}
	if n != 1 {
		// We got more than one, something is wrong
		logger.Error(vars["id"], "found ", n)
		return
	}
	data := struct {
		Item  filemaker.Record
		Items filemaker.Records
	}{}
	data.Item = reqs[0]
	//log.Println(data, n)
	params = make([]filemaker.SearchParam, 0)
	params = append(params, filemaker.SearchParam{Op: filemaker.Equal, Name: "Ordernummer", Value: data.Item["Ordernummer"].String()})
	reqs, n, err = fms.Get("G-Smart", "www_orderrader", params, 0, 30)
	if err != nil {
		logger.Error(err)
		return
	}
	data.Items = reqs
	//log.Println(data, n)

	t, _ := template.ParseFiles("html/offert.html")
	if err := t.Execute(w, data); err != nil {
		logger.Error(err)
	}
}
