package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func ReceiveUserProfileImage(w http.ResponseWriter, r *http.Request) {
	// receives file from a form in the profile page
	r.ParseMultipartForm(10 << 20) //10 MB
	file, handle, err := r.FormFile("profile-image")
	if err != nil {
		WriteInternalServerError(w, r, err.Error())
		return
	}
	defer file.Close()
	
	mimeType := handle.Header.Get("Content-Type")
	if mimeType != "image/png" {
		fmt.Println("bad file type")
		WriteInternalServerError(w, r, "incorrect file type")
		return
	}
	
	fileBuffer := make([]byte, handle.Size)
	readCount, err := file.Read(fileBuffer)
	fmt.Println(readCount)
	if err != nil {
		fmt.Println(err.Error())
		WriteInternalServerError(w, r, err.Error())
		return
	}

	id, err := GetUserIdFromCookie(r)
	if err != nil {
		fmt.Println(err.Error())
		WriteInternalServerError(w, r, err.Error())
		return
	}

	f, err := os.Create("./src/static/images/profiles/" + strconv.Itoa(id) + ".png")
	if err != nil {
		fmt.Println(err.Error())
		WriteInternalServerError(w, r, err.Error())
		return
	}
	defer f.Close()
	_, err = f.Write(fileBuffer)
	f.Sync()
	if err != nil {
		fmt.Println(err.Error())
		WriteInternalServerError(w, r, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("made it")
}

