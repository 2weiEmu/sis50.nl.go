package main

import "fmt"


var stateList = []string{"present", "absent", "cooking", "uncertain", "maybe-cooking", "cant-cook"}
var personList = []string{"rick", "youri", "robert", "milan"}
var dayList = []string{"ma", "di", "wo", "do", "vr", "za", "zo"}

type Calendar struct {
	Day [][]int
}

const DAY_COUNT = 7
const PERSON_COUNT = 4

func InitCalendarDefault() Calendar {
	var cal Calendar
	cal.Day = make([][]int, DAY_COUNT)
	for i := 0; i < DAY_COUNT; i++ {
		cal.Day[i] = make([]int, PERSON_COUNT)
		for j := range cal.Day[i] {
			cal.Day[i][j] = 0
		}
	}

	fmt.Println("[INFO] InitCalendarDefault", cal)
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
