package main

import (
	"log/slog"
	"net/http"
)

type AccessLoggerWrapper struct {
	handler http.Handler
	Logger slog.Logger
}


