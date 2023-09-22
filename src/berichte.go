package main

import "database/sql"

type BerichtQuery struct {
	Id int 
	Content string
}

func GetBerichtePrepared(preparedStatement *sql.Stmt) []BerichtQuery {

	rows, err := preparedStatement.Query()

	if err != nil {
		// TODO:
	}

	var berichte []BerichtQuery

	for rows.Next() {

		var bericht BerichtQuery
		err = rows.Scan(&bericht.Id, &bericht.Content)

		if err != nil {
			// TODO:
		}

		berichte = append(berichte, bericht)
	}

	return berichte
}
