package main

import (
	"fmt"
	"net/http"
)

type MessagesError struct {
	Message string
	Err error
}

func msgErrLog(text string, err error) MessagesError {
	errorLog.Println(text)

	return MessagesError{
		Message: text,
		Err: err,
	}
}

func (err MessagesError) Error() string {
	return fmt.Sprintf("MessagesError %s, with: %v", err.Message, err.Err)
}

func errorLogAndHttpStat(writer http.ResponseWriter, err error) {
	errorLog.Println(err)
	http.Error(
		writer, err.Error(), http.StatusInternalServerError)
}
