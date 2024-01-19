package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
NOTE:
The scheme by which the messages are saved:
// right now - all in memory, and saved to a file at the end of the day
*/

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
		writer.Write([]byte("No Page\n"))
	}
	fmt.Println(messageList)
	writer.Write(page)
}

func getMessageJson(pageNumber int) []byte {
	
	if len(messageList.Pages) <= pageNumber || pageNumber < 0 {
		return nil;
	}

	// TODO: add check
	page, _ := json.Marshal(messageList.Pages[
		len(messageList.Pages) - pageNumber - 1])
	return page
}

func POSTMessage(writer http.ResponseWriter, request *http.Request) {
	var msgPost MessagePost
	err := json.NewDecoder(request.Body).Decode(&msgPost)
	if err != nil {
		// TODO:
		fmt.Println(err)
	}

	err = addMessageToList(msgPost.Message)
	if err != nil {
		writer.Write([]byte("Failed to Add\n"))
	} else {
		writer.Write([]byte("Added\n"))
	}
	fmt.Println(messageList)
}

func addMessageToList(message string) error {
	if message == "" {
		return errors.New("Empty Message")
	}

	if len(messageList.Pages[len(messageList.Pages) - 1].Message) >= 10 {
		var msgPage MessagePage
		msgPage.Message = []string{message};
		messageList.Pages = append(messageList.Pages, msgPage)
	} else {
		messageList.Pages[len(messageList.Pages) - 1].Message =
			append(
				messageList.Pages[len(messageList.Pages)-1].Message, message)
	} 
	return nil
}
