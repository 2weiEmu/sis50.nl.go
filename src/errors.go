package main

import (
	"errors"
	"fmt"
	"net/http"
)

type LocalErr struct {
	Message string
	Err error
}

func ErrLog(text string, err error) LocalErr {
	errorLog.Println(text)

	if err == nil {
		err = errors.New(text)
	}

	return LocalErr{
		Message: text,
		Err: err,
	}
}

func (err LocalErr) Error() string {
	return fmt.Sprintf("MessagesError %s, with: %v", err.Message, err.Err)
}

func errorLogAndHttpStat(writer http.ResponseWriter, err error) {
	errorLog.Println(err)
	http.Error(
		writer, err.Error(), http.StatusInternalServerError)
}
