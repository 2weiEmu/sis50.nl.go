package src

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"sis50.nl.go/pkg/lerror"
	c "sis50.nl.go/pkg/constants"
	"sis50.nl.go/pkg/logger"
)

var AllMessagesList, _ = ReadMessages(MessageList{});

type MessagePage struct {
	Message []string `json:"messages"`
}

type MessageList struct {
	Pages []MessagePage
}

// GET pages of Messages
// POST a new message to a page, and save this
func GETMessages(writer http.ResponseWriter, request *http.Request) {
	logger.InfoLog.Println("Get Messages Request", AllMessagesList)

	AllMessagesList, err := ReadMessages(AllMessagesList)
	if err != nil {
		lerror.ErrorLogAndHttpStat(writer, err)
		return
	}

	vars := mux.Vars(request)
	pageNumber, err := strconv.Atoi(vars["pageNumber"])
	if err != nil {
		lerror.ErrorLogAndHttpStat(writer, err)
		return
	}

	page := getMessageJson(pageNumber, AllMessagesList)
	if page == nil {
		logger.InfoLog.Println("messages.go: Could not get more messages")
		http.Error(
			writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	writer.Write(page)
}

func getMessageJson(pageNumber int, messageList MessageList) []byte {
	if len(messageList.Pages) <= pageNumber || pageNumber < 0 {
		return nil;
	}

	// TODO: add check
	page, _ := json.Marshal(messageList.Pages[
		len(messageList.Pages) - pageNumber - 1])
	return page
}

func POSTMessage(writer http.ResponseWriter, request *http.Request) {
	logger.InfoLog.Println("POSTing message")
	AllMessagesList, err := ReadMessages(AllMessagesList)
	if err != nil {
		lerror.ErrorLogAndHttpStat(writer, err)
		return
	}

	var msgPost string
	err = json.NewDecoder(request.Body).Decode(&msgPost)
	if err != nil {
		lerror.ErrorLogAndHttpStat(writer, err)
		return
	}

	err = addMessageToList(msgPost, AllMessagesList)
	if err != nil {
		writer.Write([]byte("Failed to Add\n"))
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.Write([]byte("Added\n"))
	}
	logger.InfoLog.Println(AllMessagesList)
	saveMessages(AllMessagesList)
}

func addMessageToList(message string, messageList MessageList) error {
	if message == "" {
		return errors.New("Empty Message")
	}

	pagesLen := len(messageList.Pages)
	logger.InfoLog.Println(messageList.Pages)
	

	if len(messageList.Pages) == 0 {
		logger.InfoLog.Println("There are no pages...")
		messageList.Pages = []MessagePage{{Message: []string{message}}}

	} else if len(messageList.Pages[pagesLen - 1].Message) >= 10 {
		logger.InfoLog.Println("Page is full")
		var msgPage MessagePage
		msgPage.Message = []string{message};
		messageList.Pages = append(messageList.Pages, msgPage)

	} else {
		logger.InfoLog.Println("Adding to page")
		messageList.Pages[pagesLen - 1].Message =
			append(
				messageList.Pages[pagesLen-1].Message, message)

	} 
	logger.InfoLog.Println(messageList.Pages)
	return nil
}

func saveMessages(messageList MessageList) {
	logger.InfoLog.Println("writing:", messageList)
	err := os.Truncate(c.MessageFile, 0)
	if err != nil {
		lerror.ErrLog("Failed to truncate message file", err)
	}

	file, err := os.OpenFile(
		c.MessageFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		lerror.ErrLog("Failed to open message file", err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	for _, page := range messageList.Pages {
		logger.InfoLog.Println("Writing Page:", page)
		err = csvWriter.Write(page.Message)
		if err != nil {
			lerror.ErrLog("CSV writer failed when writing a message", err)
		}
	}
	csvWriter.Flush()
	err = csvWriter.Error()
	if err != nil {
		lerror.ErrLog("The CSV writer errored in some way:", err)
	}

}

func ReadMessages(messageList MessageList) (MessageList, error) {
	file, err := os.OpenFile(
		c.MessageFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		return MessageList{}, lerror.ErrLog("Could not open message file", err)
	}
	defer file.Close()

	csvr := csv.NewReader(file)
	csvr.FieldsPerRecord = -1
	records, err := csvr.ReadAll()
	if err != nil {
		return MessageList{}, lerror.ErrLog("Failed with record err", err)
	}

	messageList.Pages = []MessagePage{}
	
	for _, record := range records {
		messageList.Pages = append(messageList.Pages, MessagePage{record})
	}

	return messageList, nil
}
