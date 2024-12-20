package src

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"

	"sis50/pkg/auth"
	c "sis50/pkg/constants"
	"sis50/pkg/lerror"
	"sis50/pkg/lformatting"
	"sis50/pkg/calendar"
	"github.com/gorilla/mux"
)

type IndexPageStruct struct {
	Message string
	Args string
	ProfilePicture string
	OuterCookTable template.HTML
}

type HTMLContext struct {
	Secure string
	ConnectionLocation string
	InfoLog, RequestLog, ErrorLog *log.Logger 
	IndexTemplate *template.Template
}

func (ctx *HTMLContext) HandleIndex(w http.ResponseWriter, r *http.Request) {
	var titleMsg string
	pagesLength := len(AllMessagesList.Pages)
	if pagesLength == 0 {
		titleMsg = "No messages."
	} else {
		titleMsg = AllMessagesList.Pages[pagesLength - 1].Message[
			len(AllMessagesList.Pages[pagesLength - 1].Message) - 1]
	}

	userId, err := auth.GetUserIdFromCookie(r)
	if err != nil {
		lerror.WriteInternalServerError(w, r, err.Error())
		return
	}

	err = ctx.IndexTemplate.Execute(w, IndexPageStruct{
		OuterCookTable: template.HTML(lformatting.IndexPageTable(&calendar.StateCalendar)),
		Message: titleMsg,
		Args: ctx.ConnectionLocation + " " + ctx.Secure,
		ProfilePicture: strconv.Itoa(userId),
	})
	if err != nil {
		// TODO: move around logging
		// lerror.ErrLog("Failed to execute index template", err)
		ctx.ErrorLog.Println("Failed to execute index template")
	}
}

// could probably combine this with handlepage in a better way but rn idc, and
// yes i do parse the templates each time i need to change that
func (ctx *HTMLContext) HandleLogin(w http.ResponseWriter, r *http.Request) {
	jsArguments := ctx.ConnectionLocation + " " + ctx.Secure
	pageLocation := "src/static/templates/login.html"
	tmpl, err := template.ParseFiles(pageLocation)
	if err != nil {
		lerror.ErrLog("Could not parse template file for page", err)
		lerror.WriteInternalServerError(w, r, "Failed to parse the login template file")
		return
	}
	
	fmt.Println("Javascript Arguments before Executing:", jsArguments)
	err = tmpl.Execute(w, jsArguments)
	if err != nil {
		lerror.ErrLog("Could not execute template file", err)
		lerror.WriteInternalServerError(w, r, "Failed to execute the login template")
	}
}

func (ctx *HTMLContext) HandlePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	jsArguments := ctx.ConnectionLocation + " " + ctx.Secure

	if !slices.Contains(c.GetValidPages(), page) {
		lerror.WriteNotFound(w, r, "This doesn't seem to be a valid page")
		return
	}

	pageLocation := "src/static/templates/" + page + ".html"
	tmpl, err := template.ParseFiles(pageLocation)
	if err != nil {
		lerror.ErrLog("Could not parse template file for page", err)
		lerror.WriteInternalServerError(w, r, "Failed to parse the page template file")
		return
	}

	userId, err := auth.GetUserIdFromCookie(r)
	if err != nil {
		lerror.WriteInternalServerError(w, r, err.Error())
		return
	}
	
	fmt.Println("Javascript Arguments before Executing:", jsArguments)
	err = tmpl.Execute(w, IndexPageStruct{
		Args: jsArguments,
		ProfilePicture: strconv.Itoa(userId),
	})
	if err != nil {
		lerror.ErrLog("Could not execute template file", err)
		lerror.WriteInternalServerError(w, r, "Failed to execute the page template")
	}
}

func NewHTMLContext(
	loggerFlags int, logFile *os.File, secure string,
	connectionLocation string,
) (HTMLContext, error) {
	tmpl, err := template.ParseFiles(c.IndexFile)
	if err != nil {
		return HTMLContext{}, err
	}

	return HTMLContext{
		Secure: secure,
		ConnectionLocation: connectionLocation,
		InfoLog: log.New(logFile, "[INFO] ", loggerFlags),
		RequestLog: log.New(logFile, "[REQUEST] ", loggerFlags),
		ErrorLog: log.New(logFile, "[ERROR] ", loggerFlags),
		IndexTemplate: tmpl,
	}, nil
}

type HandleFuncWrapper struct {
	HandleFunc func(http.ResponseWriter, *http.Request)
}

func HandleFuncAsHandle(handleFunc func(http.ResponseWriter, *http.Request)) HandleFuncWrapper {
	return HandleFuncWrapper{
		HandleFunc: handleFunc,
	}
}

func (wrapper HandleFuncWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrapper.HandleFunc(w, r)
}
