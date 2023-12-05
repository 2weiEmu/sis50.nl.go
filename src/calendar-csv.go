package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)


type CalendarDay struct {
	Rick int
	Youri int
	Robert int
	Milan int
	Date time.Time
	Weekday string
}

/*
Time format of the calendar:

DD-MM-YYYY, WEEKDAY, RICK_STATE, YOURI_STATE, ROBERT_STATE, MILAN_STATE

*/

const CALENDAR_CSV = "resources/calendar.csv"

// TODO: I know loading in batches would be better, can be improved
// TODO: make it so that there is an automatic different file for each month
// TODO: perhaps adding scanning in reverse? Idk
func LoadCalendarDay(date time.Time) CalendarDay {

	loaded_day := CalendarDay {
	}

	file, err := os.Open(CALENDAR_CSV)
	if err != nil {
		// TODO:
	}
	defer file.Close()

	csv_scanner := bufio.NewScanner(file)

	year, month, day := date.Date()
	given_date := strconv.Itoa(day) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(year)
	found_date := ""
	line := ""
	var line_components []string;

	for found_date != given_date {
		csv_scanner.Scan()
		line = csv_scanner.Text()
		line_components = strings.Split(line, ",")
		found_date = line_components[0]
	}

	loaded_day.Date = date
	loaded_day.Weekday = strings.TrimSpace(line_components[1])
	loaded_day.Rick, _ = strconv.Atoi(strings.TrimSpace(line_components[2]))
	loaded_day.Youri, _ = strconv.Atoi(strings.TrimSpace(line_components[3]))
	loaded_day.Robert, _ = strconv.Atoi(strings.TrimSpace(line_components[4]))
	loaded_day.Milan, _ = strconv.Atoi(strings.TrimSpace(line_components[5]))

	fmt.Println(loaded_day)
	return loaded_day
}

func SaveCalendarDay(day CalendarDay) {

}

