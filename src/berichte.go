package main

import (
	"database/sql"
	"fmt"
)

type BerichtQuery struct {
	Id int 
	Content string
}


func GetAllBerichte(db *sql.DB) []BerichtQuery {

	var result []BerichtQuery

	rows, err := db.Query("SELECT * FROM berichte")
	defer rows.Close()
	if err != nil {
		// TODO:
	}

	for rows.Next() {
		var bericht BerichtQuery
		err := rows.Scan(&bericht.Id, &bericht.Content)
		if err != nil {
			// TODO:
		}
		result = append(result, bericht)
	}

	return result
}

func WriteNewBerichtGetId(db *sql.DB, bericht BerichtQuery) int {
	_, err := db.Exec("INSERT INTO berichte (content) VALUES (?)", bericht.Content)
	if err != nil {
		// TODO:
	}

	row, err := db.Query("SELECT id FROM berichte WHERE content = ?", bericht.Content)
	defer row.Close()
	if err != nil {
		// TODO:
	}
	row.Next()

	var result int
	row.Scan(&result)
	return result
}

func RemoveBerichtById(db *sql.DB, id int) {
	fmt.Println("Deleting with id:", id)
	_, err := db.Exec("DELETE FROM berichte WHERE id = ?", id)
	if err != nil {
		// TODO:
		fmt.Println("An error occured when deleting a bericht.")
	}
	fmt.Println("deleted...")
}
