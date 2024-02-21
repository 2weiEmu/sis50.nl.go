package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type MessagePost struct {
	Message string `json:"message"`
}

type MessagePage struct {
	Message []string `json:"messages"`
}

type MessageList struct {
	Pages []MessagePage
}

// GET pages of Messages
// POST a new message to a page, and save this
func GETMessages(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("[INFO] GET REQUEST RECEIVED")
	allMessagesList = readMessages(allMessagesList)

	vars := mux.Vars(request)
	pageNumber, err := strconv.Atoi(vars["pageNumber"])
	if err != nil {
		// TODO:
		fmt.Println(err)
	}

	// TODO: this needs more error checking
	// page, err := json.Marshal(messageList.Pages[pageNumber])
	page := getMessageJson(pageNumber)

	if page == nil {
		// TODO:
		fmt.Println(err)
		http.Error(
			writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	fmt.Println(allMessagesList)
	writer.Header().Set("Access-Control-Allow-Headers", "x-requested-with")
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
	fmt.Println("[INFO] POST REQUEST RECEIVED")
	allMessagesList = readMessages(allMessagesList)
	var msgPost MessagePost
	err := json.NewDecoder(request.Body).Decode(&msgPost)
	if err != nil {
		// TODO:
		fmt.Println(err)
	}

	fmt.Println(allMessagesList)
	err = addMessageToList(msgPost.Message)
	if err != nil {
		writer.Write([]byte("Failed to Add\n"))
	} else {
		writer.Write([]byte("Added\n"))
	}
	fmt.Println(allMessagesList)
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
		// TODO:
		fmt.Println(err)
	}

	file, err := os.OpenFile(
		MessageFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		// TODO:
		fmt.Println(err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	for _, page := range messageList.Pages {
		fmt.Println("[INFO] Writing Page:", page)
		err = csvWriter.Write(page.Message)
		if err != nil {
			// TODO:
			fmt.Println(err)
		}
	}
	csvWriter.Flush()
	err = csvWriter.Error()
	if err != nil {
		// TODO:
		fmt.Println(err)
	}

}

func readMessages(messageList MessageList) MessageList {

	file, err := os.OpenFile(
		MessageFile, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		// TODO:
		fmt.Println(err)
	}
	defer file.Close()

	csvr := csv.NewReader(file)
	csvr.FieldsPerRecord = -1
	records, err := csvr.ReadAll()
	if err != nil {
		// TODO:
		fmt.Println("[ERROR] Record:", err)
	}

	messageList.Pages = []MessagePage{}
	
	for _, record := range records {
		// fmt.Println("[INFO] Read record:", record)
		messageList.Pages = append(messageList.Pages, MessagePage{record})
	}
	// fmt.Println("[INFO] Finished Read messageList:", messageList)

	return messageList
}
