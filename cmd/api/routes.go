package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Endpoints goes here
	mux.HandleFunc("POST /v1/customers", app.registerCustomerHandler)

	return mux
}
