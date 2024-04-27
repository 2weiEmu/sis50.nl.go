package main

import (
	"log"
	"log/slog"
	"net/http"
)

const loggerFlags = log.LstdFlags | log.Llongfile | log.Ldate | log.Ltime

type AccessLoggerWrapper struct {
	handler http.Handler
	Logger slog.Logger
}
