export function getRelevantTableElement(week, person, day) {
	var week_table = document.getElementsByClassName(week + "-week")[0]
	var element = week_table.getElementsByClassName(day)[0].getElementsByClassName(person)[0]

	return element
}

export function send_day_toggle(socket_conn, event, week) {
	// var target = event.target
	var person = event.target.className
	var day = event.target.parentNode.className


	if (event.button != 0) { return }

	if (people.indexOf(person) == -1 || days.indexOf(day) == -1) {
		return


	} else {
		console.log(person, day)
		socket_conn.send(JSON.stringify({
				command: "toggle",
				currentState: "empty",
				week: week,
				person: person,
				day: day
			})
		)
	}

}

export var current_week_table = document.getElementsByClassName("current-week")[0];
export var next_week_table = document.getElementsByClassName("next-week")[0];
export var people = ["rick", "youri", "robert", "milan"]
export var days = ["mandaag", "dinsdag", "woensdag", "dondersdag", "vrijdag", "zaterdag", "zondag"]
