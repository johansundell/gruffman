package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	routes = append(routes, Route{"offert", "GET", "/offert/{id}", offertHandler})
}

func offertHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprint(w, vars["id"])
}
