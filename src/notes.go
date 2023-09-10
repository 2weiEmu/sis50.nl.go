package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Note struct {
	Content string;
	Week string;
	Day string;
	Person string;
}

func loadNotes() []Note {
	file, err := os.Open("./src/resources/notes")

	if err != nil {
		log.Fatal("Failed to load notes with error", err)
	}

	defer file.Close()

	var result []Note;

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		note := ParseNote(line)

		result = append(result, note)
	}

	fmt.Println(result)

	return result
}

func saveNotes(notes []Note) {

	err := os.Truncate(".src/resource/notes", 0)

	if err != nil {
		// TODO:
	}

	file, err := os.OpenFile("./src/resources/notes", os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println("failed to open save file...")
	}

	defer file.Close()

	for _, note := range notes {

		_, err := file.WriteString(NoteString(note))

		if err != nil {
			fmt.Println("failed writling line with error", err)

		}
	}

	fmt.Println("finished saving notes to file")
}

func NoteString(note Note) string {
	return "_$" + note.Content + "$" + note.Week +"$" + note.Person + "$" + note.Day
}

func ParseNote(message string) Note {

	// note message: notecmd$content$week$person$day 
	// TODO: make everything Json based I guess, so that content can include dollar signs

	arr := strings.Split(message, "$")

	return Note {
		Content: arr[1],
		Week: arr[2],
		Day: arr[4],
		Person: arr[3],
	}
}
