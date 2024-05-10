package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type HTMLContext struct {
	Secure string
	ConnectionLocation string
	InfoLog, RequestLog, ErrorLog *log.Logger 
	IndexTemplate *template.Template
}

func (ctx *HTMLContext) HandleIndex(w http.ResponseWriter, r *http.Request) {
	var titleMsg string
	pagesLength := len(allMessagesList.Pages)
	if pagesLength == 0 {
		titleMsg = "No messages."
	} else {
		titleMsg = allMessagesList.Pages[pagesLength - 1].Message[
			len(allMessagesList.Pages[pagesLength - 1].Message) - 1]
	}

	err := ctx.IndexTemplate.Execute(w, IndexPageStruct{
		Message: titleMsg,
		Args: ctx.ConnectionLocation + " " + ctx.Secure,
	})
	if err != nil {
		// TODO: move around logging
		// ErrLog("Failed to execute index template", err)
		ctx.ErrorLog.Println("Failed to execute index template")
	}
}

func (ctx *HTMLContext) HandlePage(w http.ResponseWriter, r *http.Request) {

}


func initHTMLContext(
	loggerFlags int, logFile *os.File, secure string,
	connectionLocation string,
) (HTMLContext, error) {
	tmpl, err := template.ParseFiles(IndexFile)
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

