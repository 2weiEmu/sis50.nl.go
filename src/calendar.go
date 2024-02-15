package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

type CalMessage struct {
	Day string `json:"day"`
	Person string `json:"person"`
	State string `json:"state"`
}

type Calendar struct {
	Day [][]int
}

func InitCalendarDefault() Calendar {
	var newCal Calendar
	newCal.Day = make([][]int, DayCount)
	for i := 0; i < DayCount; i++ {
		newCal.Day[i] = make([]int, UserCount)
		for j := range newCal.Day[i] {
			newCal.Day[i][j] = 0
		}
	}

	fmt.Println("[INFO] InitCalendarDefault", newCal)
	return newCal
}

func ReadCalendar(cal Calendar) Calendar {
	
	calFile, err := os.Open(CalendarFile)
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
		for i := 0; i < UserCount; i++ {
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

func UpdateCalendar(cal Calendar, message CalMessage) string {

	var dayIndex, personIndex int

	for i, d := range DayList {
		if d == message.Day {
			dayIndex = i 
			break
		}
	}

	for i, p := range ConstPersonList {
		if p == message.Person {
			personIndex = i 
			break
		}
	}

	new_state := (cal.Day[dayIndex][personIndex] + 1) % len(StateList)
	cal.Day[dayIndex][personIndex] = new_state

	return StateList[new_state]
}

func WriteCalendar(cal Calendar) {
	fmt.Println("Saving Calendar")
	err := os.Truncate(CalendarFile, 0)
	if err != nil {
		// TODO:
		fmt.Println("[ERROR]", err)
	}

	file, err := os.OpenFile(CalendarFile, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		// TODO:
		fmt.Println("[ERROR]", err)
	}
	defer file.Close()

	for _, r := range cal.Day {
		f := ""
		for _, c := range r {
			f += strconv.Itoa(c)
		}
		f += "\n"
		_, err := file.WriteString(f)
		if err != nil {
			fmt.Println("[ERROR] Writing to File", err)
		}
	}
}

func DayWebsocketHandler(conn *websocket.Conn) {
	fmt.Println("Activating WebSocket handler...")

	WebSocketDayConnections = append(WebSocketDayConnections, conn)
	fmt.Println(WebSocketDayConnections)

	var message CalMessage
	for {
		err := websocket.JSON.Receive(conn, &message)
		if err != nil {
			// TODO:
			fmt.Println(err)
			break
		}
		fmt.Println("Message received: ", message)

		if message.State != "open-calendar" {
			message.State = UpdateCalendar(StateCalendar, message)
			BroadcastToConnections(message)
		} else {
			m := ""
			for _, s := range StateCalendar.Day {
				for _, k := range s {
					m += strconv.Itoa(k)
				}
				m += "/"
			}
			fmt.Println("[INFO] Open:", m)

			message.Day = m
			err := websocket.JSON.Send(conn, &message)
			if err != nil {
				fmt.Println(err)
				// TODO:
			}
		}
	}
	WriteCalendar(StateCalendar)
	WebSocketDayConnections = RemoveWebsocketFromPool(conn, WebSocketDayConnections)
}
