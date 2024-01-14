package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)


var stateList = []string{"present", "absent", "cooking", "uncertain", "maybe-cooking", "cant-cook"}
var personList = []string{"rick", "youri", "robert", "milan"}
var dayList = []string{"ma", "di", "wo", "do", "vr", "za", "zo"}

type Calendar struct {
	Day [][]int
}

const DAY_COUNT = 7
const PERSON_COUNT = 4

const CALENDAR_FILE = "./resources/calendar"

func InitCalendarDefault() Calendar {
	var newCal Calendar
	newCal.Day = make([][]int, DAY_COUNT)
	for i := 0; i < DAY_COUNT; i++ {
		newCal.Day[i] = make([]int, PERSON_COUNT)
		for j := range newCal.Day[i] {
			newCal.Day[i][j] = 0
		}
	}

	fmt.Println("[INFO] InitCalendarDefault", newCal)
	return newCal
}

func ReadCalendar(cal Calendar) Calendar {
	
	calFile, err := os.Open(CALENDAR_FILE)
	if err != nil {
		// TODO:
		fmt.Println(err)
	}
	defer calFile.Close()

	fileScanner := bufio.NewScanner(calFile)
	fileScanner.Split(bufio.ScanLines)

	d := 0
	for fileScanner.Scan() {
		tState := strings.Split(fileScanner.Text(), "")
		for i := 0; i < PERSON_COUNT; i++ {
			cal.Day[d][i], err = strconv.Atoi(tState[i])
			if err != nil {
				// TODO:
				fmt.Println(err)
			}
		}
		d++
	}

	fmt.Println("[INFO] Calendar Read:", cal)

	return cal
}

func UpdateCalendar(cal Calendar, message MessageStruct) string {

	var dayIndex, personIndex int

	for i, d := range dayList {
		if d == message.Day {
			dayIndex = i 
			break
		}
	}

	for i, p := range personList {
		if p == message.Person {
			personIndex = i 
			break
		}
	}

	new_state := (cal.Day[dayIndex][personIndex] + 1) % len(stateList)
	cal.Day[dayIndex][personIndex] = new_state

	return stateList[new_state]
}

func WriteCalendar(cal Calendar) {
}
