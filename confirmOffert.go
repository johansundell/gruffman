package main

import "net/http"

func init() {
	routes = append(routes, Route{"confirmOffert", "GET", "/confirmoffert/{id}", confirmOffert})
}

func confirmOffert(w http.ResponseWriter, r *http.Request) {
	// TODO: Write this function ;)
}
