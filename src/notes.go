package main

import (
	"database/sql"
	"fmt"
)

type Note struct {
	Content string;
	Week string;
	Day string;
	Person string;
	Id string
}

func loadNotes(db *sql.DB) []Note {
    rows, err := db.Query(`SELECT * FROM notes`)
    defer rows.Close()

    if err != nil {
        fmt.Println("Prepared statement failed to execute with error:", err)
		return nil
    }
    var result []Note;

    for rows.Next() {
        var note Note;
        err = rows.Scan(&note.Id, &note.Week, &note.Person, &note.Day, &note.Content)

        if err != nil {
            fmt.Println("Failed to retrieve row:", rows, "With the error:", err)
			return nil
        }
        result = append(result, note)
    }
    return result
}

func removeNoteById(db *sql.DB, id int) {
	_, err := db.Exec(`DELETE FROM notes WHERE id = ?`, id)
	if err != nil {
		// TODO:
		fmt.Println("Error here (3)...", err)
	}
}

func saveNoteGetID(db *sql.DB, note Note) int {
	_, err := db.Exec(`INSERT INTO notes (week, person, day, content) VALUES (?, ?, ?, ?)`, note.Week, note.Person, note.Day, note.Content)
	if err != nil {
		// TODO:
		fmt.Println("Error here (1)...", err)
	}

	// this is by far not optimal
	row, err := db.Query(`SELECT id FROM notes WHERE week = ? AND person = ? AND day = ? AND content = ?`, note.Week, note.Person, note.Day, note.Content)
	defer row.Close()
	if err != nil {
		// TODO:
		fmt.Println("Error here (2)...", err)
	}
	row.Next()

	var result int
	row.Scan(&result)
	return result

}
