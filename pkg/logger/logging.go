package logger

import (
	"log"
)

var InfoLog, RequestLog, ErrorLog *log.Logger

const LoggerFlags = log.LstdFlags | log.Lshortfile | log.Ldate | log.Ltime

type AccessLoggerWrapper struct {
	infoLog, requestLog, errorLog *log.Logger
}
