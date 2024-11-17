package calendar

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
	"sis50.nl.go/pkg/logger"
	n "sis50.nl.go/pkg/node"
	c "sis50.nl.go/pkg/constants"
	"sis50.nl.go/pkg/lerror"
)

var StateCalendar = ReadCalendar(InitCalendarDefault())

type CalendarHandler struct {
	Connections []*websocket.Conn
	// TODO: need to find some way to include loggers for everyone
	InfoLog, RequestLog, ErrorLog *log.Logger 
}

func NewCalendarHandler(
	loggerFlags int, logFile *os.File,
) CalendarHandler {
	return CalendarHandler { 
		Connections: make([]*websocket.Conn, 0),
		InfoLog: log.New(logFile, "[INFO] ", loggerFlags),
		RequestLog: log.New(logFile, "[REQUEST] ", loggerFlags),
		ErrorLog: log.New(logFile, "[ERROR] ", loggerFlags),
	}
}

func (handler *CalendarHandler) HandleCalendarWebsocket(conn *websocket.Conn) {
	handler.InfoLog.Println("Activating WebSocket handler...")

	handler.Connections = append(handler.Connections, conn)
	var message CalMessage
	for {
		err := websocket.JSON.Receive(conn, &message)
		if err != nil {
			// ErrLog("Failed to read the websocket message", err)
			break
		}

		if message.State == "OPEN" {

		} else {
			message.State = UpdateCalendar(StateCalendar, message)
			handler.BroadcastToConnections(message)
		}
	}
	WriteCalendar(StateCalendar)
	handler.Connections = n.RemoveWebsocketFromPool(conn, handler.Connections)
}

type CalMessage struct {
	Day string `json:"day"`
	Person string `json:"person"`
	State string `json:"state"`
}

type Calendar struct {
	Day [][]int
}

var calFileReady = make(chan bool)

func (handler *CalendarHandler) ShiftCalendarDaily() {
	for {
		t := time.Now().AddDate(0, 0, 1)
		targetTime := time.Date(
			t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Now().Location())
		timeUntilTmr := time.Until(targetTime)
		time.Sleep(timeUntilTmr)

		WriteCalendar(StateCalendar)
		StateCalendar = shiftCalendar()
		WriteCalendar(StateCalendar)

		handler.BroadcastToConnections(genOpenCalMessage())
	}
}

func shiftCalendar() Calendar {
	cal := ReadCalendar(InitCalendarDefault())

	for d := 0; d < len(c.GetDayList()) - 1; d++ {
		for s := 0; s < len(c.GetPersonList()); s++ {
			cal.Day[d][s] = cal.Day[d+1][s];
		}
	}

	lastDay := len(c.GetDayList()) - 1
	for s := 0; s < len(c.GetPersonList()); s++ {
		cal.Day[lastDay][s] = 0;
	}

	return cal
}

func InitCalendarDefault() Calendar {
	var newCal Calendar
	newCal.Day = make([][]int, c.DayCount)
	for i := 0; i < c.DayCount; i++ {
		newCal.Day[i] = make([]int, c.UserCount)
		for j := range newCal.Day[i] {
			newCal.Day[i][j] = 0
		}
	}

	return newCal
}

func ReadCalendar(cal Calendar) Calendar {
	
	calFile, err := os.Open(c.CalendarFile)
	if err != nil {
		lerror.ErrLog("Failed to open the calendar file", err)
	}
	defer calFile.Close()

	fileScanner := bufio.NewScanner(calFile)
	fileScanner.Split(bufio.ScanLines)

	d := 0
	for fileScanner.Scan() {
		tState := strings.Split(fileScanner.Text(), "")
		for i := 0; i < c.UserCount; i++ {
			cal.Day[d][i], err = strconv.Atoi(tState[i])
			if err != nil {
				lerror.ErrLog("Failed to convert state string when reading calendar", err)
			}
		}
		d++
	}

	return cal
}

func IndexInStateList(state string) int {
	for i, e := range c.GetStateList() {
		fmt.Println(e, state)
		if e == state {
			return i
		}
	}

	return 0
}

func UpdateCalendar(cal Calendar, message CalMessage) string {

	var dayIndex, personIndex int

	for i, d := range c.GetDayList() {
		if d == message.Day {
			dayIndex = i 
			break
		}
	}

	for i, p := range c.GetPersonList() {
		if p == message.Person {
			personIndex = i 
			break
		}
	}

	var new_state int
	setMsg := strings.Split(message.State, ":")
	if setMsg[0] == "SET" {
		new_state = IndexInStateList(string(setMsg[1]))
		fmt.Println(setMsg, new_state)
		cal.Day[dayIndex][personIndex] = new_state
	} else {
		new_state = (cal.Day[dayIndex][personIndex] + 1) % len(c.GetStateList())
	}

	cal.Day[dayIndex][personIndex] = new_state

	return c.GetStateList()[new_state]
}

func WriteCalendar(cal Calendar) {
	logger.InfoLog.Println("Saving Calendar")
	err := os.Truncate(c.CalendarFile, 0)
	if err != nil {
		lerror.ErrLog("Failed to truncate file // TODO: THIS SHOULD BE REMOVED", err)
	}

	file, err := os.OpenFile(c.CalendarFile, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		lerror.ErrLog("Failed to open calendar file", err)
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
			lerror.ErrLog("Something went wrong writing to the calendar file", err)
		}
	}
}

func resetCalendar() {
	StateCalendar = InitCalendarDefault()
	WriteCalendar(StateCalendar);
}

func genOpenCalMessage() CalMessage {
	m := ""
	for _, s := range StateCalendar.Day {
		for _, k := range s {
			m += strconv.Itoa(k)
		}
		m += "/"
	}

	message := CalMessage{}
	message.Day = m
	message.State = "open-calendar" 

	return message
}

func (handler *CalendarHandler) weeklyResetTimer() {
	for {
		// wait until Monday
		currentWeekday := time.Now().Weekday()
		var targetTime time.Time
		t := time.Now().AddDate(0, 0, 7 - int(currentWeekday - 1))
		targetTime = time.Date(
			t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Now().Location())
		timeUntilMonday := time.Until(targetTime)
		time.Sleep(timeUntilMonday)

		resetCalendar()
		handler.BroadcastToConnections(genOpenCalMessage())
	}
}

func (handler *CalendarHandler) BroadcastToConnections(message CalMessage) {
	handler.InfoLog.Println("Broadcasting Websocket")
	for i := 0; i < len(handler.Connections); i++ {
		logger.InfoLog.Println("Sending to: ", handler.Connections[i])
		err := websocket.JSON.Send(handler.Connections[i], message)
		if err != nil {
			lerror.ErrLog("Failed to send JSON via websocket during broadcast", err)
		}
	}
}
