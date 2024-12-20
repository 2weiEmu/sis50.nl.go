package lerror

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"runtime"

	"sis50/pkg/logger"
)

type LocalErr struct {
	Message string
	Err error
}

func ErrLog(text string, err error) LocalErr {
	_, file, no, _ := runtime.Caller(1)
	logger.ErrorLog.Println("[IN ", file, ":", no, "]", text)

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

func ErrorLogAndHttpStat(writer http.ResponseWriter, err error) {
	logger.ErrorLog.Println(err)
	http.Error(
		writer, err.Error(), http.StatusInternalServerError)
}

func WriteInternalServerError(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl, err := template.ParseFiles("src/static/templates/500.html")
	if err != nil {
		// TODO: as we shouldnt parse the template on each go	
	}
	err = tmpl.Execute(w, message)
	if err != nil {
		// TODO: idfk
	}
}

func WriteUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	tmpl, err := template.ParseFiles("src/static/templates/401.html")
	if err != nil {
		// TODO: as we shouldnt parse the template on each go	
	}
	err = tmpl.Execute(w, message)
	if err != nil {
		// TODO: idfk
	}
}

func WriteNotFound(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusNotFound)
	tmpl, err := template.ParseFiles("src/static/templates/404.html")
	if err != nil {
		// TODO: as we shouldnt parse the template on each go	
	}
	err = tmpl.Execute(w, message)
	if err != nil {
		// TODO: idfk
	}
}
