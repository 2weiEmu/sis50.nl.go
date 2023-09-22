import * as helper from "./koken-helper.js" 

// variables
var current_week_table = document.getElementsByClassName("current-week")[0];
var next_week_table = document.getElementsByClassName("next-week")[0];

var people = ["rick", "youri", "robert", "milan"]
var days = ["mandaag", "dinsdag", "woensdag", "dondersdag", "vrijdag", "zaterdag", "zondag"]

var socket_conn = new WebSocket("ws://localhost:8000/koken-ws")

var custom_context = document.getElementById("custom-context-menu")
var add_note_button = document.getElementById("add-note-button")
var delete_note_button = document.getElementById("remove-note-button")

// websocket handling
socket_conn.onopen = function() {
	socket_conn.send(JSON.stringify({
		command: "open"
	}))
}

socket_conn.onmessage = function(event) {
	console.log("event.data:", event.data)

	var message = JSON.parse(event.data)
	var command = message.command

	if (command == "toggle") {

		if (message.currentState == "E" ) { 
			message.currentState = "_" 
		}

		var element = helper.getRelevantTableElement(message.week, message.person, message.day)
		element.innerHTML = message.currentState + element.innerHTML.slice(1)


	} else if (command == "addnote") {

		var element = helper.getRelevantTableElement(message.week, message.person, message.day)

		// TODO: this can definitely be improved
		// we have definitely reached the point where an ID would be easier
		// we need to safely construct this message using javascript elements, instead of just a raw string
		element.innerHTML += `<div class="note"> 
								<div class="main-note" name="${message.week}$${message.day}$${message.person}" onclick="revealNote(this)">
									<div class="note-container" style="display: none;">
										<div class="note-content" name="${message.week}$${message.day}$${message.person}">${message.currentState}</div>
									</div>
								</div>
								<div class="note-close-button" style="display: none;" onclick="closeNote(this)">X</div>
							  </div>`
		

	} else if (command == "deletenote") {
		var arr = message.split("$")	

		var content, week, person, day

		[_, content, week, person, day] = arr

		var element = helper.getRelevantTableElement(week, person, day)

		var localNoteList = element.children.getElementsByClassName("note-content")

		var removeNode

		for (var i = 0; i < localNoteList.length; i++) {
			if (localNoteList[i].innerHTML == content) {
				removeNode = localNoteList[i]
				break
			}
		}

		element.removeChild(removeNode)
	}

}

function get_command(message) {
	return message.split("$")[0]
}

// day toggling mechanisms
current_week_table.addEventListener("mousedown", function(ev) {
	send_day_toggle(ev, "current")
})

next_week_table.addEventListener("mousedown", function(ev) {
	send_day_toggle(ev, "next")
})

function send_day_toggle(event, week) {

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


// key bindings
document.addEventListener("keyup", function(ev) {
	if (ev.key == "Escape") {
		console.log("escape key pressed")
		custom_context.style.display = "none"
	}
})

// context menu
current_week_table.addEventListener("contextmenu", function(ev) { // TODO: make the same for next week
	ev.preventDefault()

	console.log(ev.clientX);
	console.log(ev.clientY);

	custom_context.style.top = ev.clientY + "px"
	custom_context.style.left = ev.clientX + "px"
	custom_context.style.display = "inline-block";
	
	var person = ev.target.className
	var day = ev.target.parentNode.className

	var delete_name = ev.target.name

	console.log(person, day)
	add_note_button.setAttribute("name", "current$" + person + "$" + day)
	delete_note_button.setAttribute("name", delete_name)

})

// admin panel
var admin_panel_toggle = document.getElementById("admin-panel-toggle")
admin_panel_toggle.addEventListener("click", toggle_admin_panel)

function toggle_admin_panel() {

	var display = admin_panel_toggle.parentNode.children[1].style.display
	admin_panel_toggle.parentNode.children[1].style.display = display == "none" ? "flex" : "none"
	console.log("toggled admin panel")
}

// berichte
var new_bericht_field = document.getElementById("new-bericht-field")
var post_bericht_button = document.getElementById("post-bericht-button")

post_bericht_button.addEventListener("click", postBericht)

function postBericht() {
	var bericht = new_bericht_field.value
	new_bericht_field.value = ""

	socket_conn.send(JSON.stringify({
		command: "post-bericht",
		currentState: bericht
	}))
}
