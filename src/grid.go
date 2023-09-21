package main

import (
	"os"
	"bufio"
	"log"
	"fmt"
	"strings"
)

func loadGrid() [][]string {
	file, err := os.Open("./src/resources/grid")

	if err != nil {
		log.Fatal("failed to open grid file with error", err)
	}

	defer file.Close()

	var grid [][]string;

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		fmt.Println(line)

		grid = append(grid, line)
	}
	return grid;
}

func ModifyGrid(grid [][]string, week, person, day string) string {

	row := 0
	if week == "next" {
		row += 7
	}
	row += findIndex(dayList, day)
	col := findIndex(personList, person)

	grid[row][col] = GetNextMark(grid[row][col])

	return grid[row][col]
}

func GetNextMark(sy string) string {
	return gridStates[(findIndex(gridStates, sy) + 1) % len(gridStates)]
}

func saveGrid(grid [][]string) {
	_ = os.Truncate("./src/resources/grid", 0)

	file, err := os.OpenFile("./src/resources/grid", os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println("failed to open save file")
	}

	defer file.Close()

	for _, arr := range grid {
		_, err := file.WriteString(strings.Join(arr, ",") + "\n")

		if err != nil {
			log.Fatal("failed writing line with error", err)
		}
	}

	fmt.Println("Finished saving grid to file.")

}
