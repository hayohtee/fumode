package main

import "github.com/hayohtee/fumode/internal/jsonlog"

// application holds the dependencies for the handlers, middlewares
// and helpers.
type application struct {
	config config
	logger *jsonlog.Logger
}
