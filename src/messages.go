package src

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/2weiEmu/sis50.nl.go/pkg/lerror"
	"github.com/2weiEmu/sis50.nl.go/pkg/logger"
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
	logger.InfoLog.Println("Get Messages Request")

	_, err := ReadMessages(AllMessagesList)
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

	page := getMessageJson(pageNumber)
	if page == nil {
		logger.InfoLog.Println("messages.go: Could not get more messages")
		http.Error(
			writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	writer.Write(page)
}

func getMessageJson(pageNumber int) []byte {
	if len(AllMessagesList.Pages) <= pageNumber || pageNumber < 0 {
		return nil;
	}

	// TODO: add check
	page, _ := json.Marshal(AllMessagesList.Pages[
		len(AllMessagesList.Pages) - pageNumber - 1])
	return page
}

func POSTMessage(writer http.ResponseWriter, request *http.Request) {
	logger.InfoLog.Println("POSTing message")
	allMessagesList, err := ReadMessages(AllMessagesList)
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

	if len(AllMessagesList.Pages) == 0 {
		AllMessagesList.Pages = []MessagePage{{Message: []string{message}}}
	} else if len(AllMessagesList.Pages[len(AllMessagesList.Pages) - 1].Message) >= 10 {
		var msgPage MessagePage
		msgPage.Message = []string{message};
		AllMessagesList.Pages = append(AllMessagesList.Pages, msgPage)
	} else {
		AllMessagesList.Pages[len(AllMessagesList.Pages) - 1].Message =
			append(
				AllMessagesList.Pages[len(AllMessagesList.Pages)-1].Message, message)
	} 
	return nil
}

func saveMessages(messageList MessageList) {
	err := os.Truncate(MessageFile, 0)
	if err != nil {
		lerror.ErrLog("Failed to truncate message file", err)
	}

	file, err := os.OpenFile(
		MessageFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
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
		MessageFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
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
