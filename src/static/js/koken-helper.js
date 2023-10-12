export function getRelevantTableElement(week, person, day) {
	var week_table = document.getElementsByClassName(week + "-week")[0]
	var element = week_table.getElementsByClassName(day)[0].getElementsByClassName(person)[0]

	return element
}

export function present_note_adding(window, event_target) {

	var rect = event_target.getBoundingClientRect()

	var left = rect.left + window.scrollX
	var top = rect.top + window.scrollY

	console.log(top, left)

	return prompt("What should the note say:")
}

export function send_day_toggle(window, socket_conn, event, week) {
	// var target = event.target
	var person = event.target.className
	var day = event.target.parentNode.className



	if (event.button != 0) { return }

	if (days.indexOf(day) != -1 && person == "dag") {

		// present note adding interface
		var note_content = present_note_adding(window, event.target)

		if (note_content.trim() != "") {
			socket_conn.send(JSON.stringify({
				command: "addnote",
				currentState: note_content,
				week: week,
				person: person,
				day: day
			})
			)
		}

		// if entry non-empty add note via ws

	} else if (people.indexOf(person) == -1 || days.indexOf(day) == -1) {
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
