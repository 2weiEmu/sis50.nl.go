package lformatting

import (
	"fmt"
	"strconv"
	"strings"

	"sis50.nl.go/pkg/calendar"
	"sis50.nl.go/pkg/constants"
)

func IndexPageTable(cal *calendar.Calendar) string {

	start := `<div class="outer-cook-table"> 
		<div class="header"><p>Rick</p></div>
		<div class="header"><p>Youri</p></div>
		<div class="header"><p>Robert</p></div>
		<div class="header"><p>Jop</p></div>
		<div class="header"><p>Dag</p></div>`
	for i, d := range constants.GetDayList() {
		var bg string = "bg"
		if i % 2 != 0 {
			bg = ""
		}
		var pers string

		for personIndex, p := range constants.GetPersonList() {
			state := constants.GetStateList()[cal.Day[i][personIndex]]
			imageSource := "/images/" + state + ".svg"

			start += fmt.Sprintf(`<div class="%s %s %s" data-person="%s" data-day="%s">
				<img data-state="%s" src="%s"/></div>
			`, p, d, bg, p, d, state, imageSource)
			pers = p
		}

		dayCount := "day" + strconv.Itoa(i + 1)
		dayCapitalised := strings.Title(d) + "."
		start += fmt.Sprintf(`
			<div class="%s %s %s %s day"><p>%s</p></div>
		`, pers, d, bg, dayCount, dayCapitalised)
	}
	start += "</div>"

	return start
}
