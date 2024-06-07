package src

import (
	"log"
)

var InfoLog, RequestLog, ErrorLog *log.Logger

const LoggerFlags = log.LstdFlags | log.Llongfile | log.Ldate | log.Ltime

type AccessLoggerWrapper struct {
	infoLog, requestLog, errorLog *log.Logger
}
