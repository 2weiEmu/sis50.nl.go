package main

import (
	"log"
)

const loggerFlags = log.LstdFlags | log.Llongfile | log.Ldate | log.Ltime

type AccessLoggerWrapper struct {
	infoLog, requestLog, errorLog *log.Logger
}
