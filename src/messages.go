package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type MessagePage struct {
	Message []string `json:"messages"`
}

type MessageList struct {
	Pages []MessagePage
}

// GET pages of Messages
// POST a new message to a page, and save this
func GETMessages(writer http.ResponseWriter, request *http.Request) {
	infoLog.Println("Get Messages Request")

	_, err := readMessages(allMessagesList)
	if err != nil {
		errorLogAndHttpStat(writer, err)
		return
	}

	vars := mux.Vars(request)
	pageNumber, err := strconv.Atoi(vars["pageNumber"])
	if err != nil {
		errorLogAndHttpStat(writer, err)
		return
	}

	page := getMessageJson(pageNumber)
	if page == nil {
		infoLog.Println("messages.go: Could not get more messages")
		http.Error(
			writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	writer.Write(page)
}

func getMessageJson(pageNumber int) []byte {
	if len(allMessagesList.Pages) <= pageNumber || pageNumber < 0 {
		return nil;
	}

	// TODO: add check
	page, _ := json.Marshal(allMessagesList.Pages[
		len(allMessagesList.Pages) - pageNumber - 1])
	return page
}

func POSTMessage(writer http.ResponseWriter, request *http.Request) {
	infoLog.Println("POSTing message")
	allMessagesList, err := readMessages(allMessagesList)
	if err != nil {
		errorLogAndHttpStat(writer, err)
		return
	}

	var msgPost string
	err = json.NewDecoder(request.Body).Decode(&msgPost)
	if err != nil {
		errorLogAndHttpStat(writer, err)
		return
	}

	err = addMessageToList(msgPost)
	if err != nil {
		writer.Write([]byte("Failed to Add\n"))
	} else {
		writer.Write([]byte("Added\n"))
	}
	saveMessages(allMessagesList)
}

func addMessageToList(message string) error {
	if message == "" {
		return errors.New("Empty Message")
	}

	if len(allMessagesList.Pages) == 0 {
		allMessagesList.Pages = []MessagePage{{Message: []string{message}}}
	} else if len(allMessagesList.Pages[len(allMessagesList.Pages) - 1].Message) >= 10 {
		var msgPage MessagePage
		msgPage.Message = []string{message};
		allMessagesList.Pages = append(allMessagesList.Pages, msgPage)
	} else {
		allMessagesList.Pages[len(allMessagesList.Pages) - 1].Message =
			append(
				allMessagesList.Pages[len(allMessagesList.Pages)-1].Message, message)
	} 
	return nil
}

func saveMessages(messageList MessageList) {
	err := os.Truncate(MessageFile, 0)
	if err != nil {
		ErrLog("Failed to truncate message file", err)
	}

	file, err := os.OpenFile(
		MessageFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		ErrLog("Failed to open message file", err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	for _, page := range messageList.Pages {
		infoLog.Println("Writing Page:", page)
		err = csvWriter.Write(page.Message)
		if err != nil {
			ErrLog("CSV writer failed when writing a message", err)
		}
	}
	csvWriter.Flush()
	err = csvWriter.Error()
	if err != nil {
		ErrLog("The CSV writer errored in some way:", err)
	}

}

func readMessages(messageList MessageList) (MessageList, error) {
	file, err := os.OpenFile(
		MessageFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		return MessageList{}, ErrLog("Could not open message file", err)
	}
	defer file.Close()

	csvr := csv.NewReader(file)
	csvr.FieldsPerRecord = -1
	records, err := csvr.ReadAll()
	if err != nil {
		return MessageList{}, ErrLog("Failed with record err", err)
	}

	messageList.Pages = []MessagePage{}
	
	for _, record := range records {
		messageList.Pages = append(messageList.Pages, MessagePage{record})
	}

	return messageList, nil
}
