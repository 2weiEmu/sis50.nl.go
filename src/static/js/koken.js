var current_week_table = document.getElementsByClassName("current-week")[0];
var next_week_table = document.getElementsByClassName("next-week")[0];

var people = ["rick", "youri", "robert", "milan"]
var days = ["mandaag", "dinsdag", "woensdag", "dondersdag", "vrijdag", "zaterdag", "zondag"]

var socket_conn = new WebSocket("ws://localhost:8000/koken-ws")

var custom_context = document.getElementById("custom-context-menu")
var add_note_button = document.getElementById("add-note-button")
var delete_note_button = document.getElementById("remove-note-button")

// websocket handling
socket_conn.onopen = function(event) {
	socket_conn.send("open$")
}

socket_conn.onmessage = function(event) {
	console.log("event.data:", event.data)
	var message = event.data

	var command = get_command(message)

	if (command == "toggle") {
		var arr = message.split("$")

		var state = arr[1]
		var e, week, person, day;
		
		[e, state, week, person, day] = arr

		console.log(state, week, person, day)

		if (state == "E") { state = "_" }

		var element = getRelevantTableElement(week, person, day)

		// for now we jankily replace the first character, as that is the only thing that should have to change
		element.innerHTML = state + element.innerHTML.slice(1)

	} else if (command == "addnote") {
		var arr = message.split("$")	

		var content, week, person, day

		[_, content, week, person, day] = arr

		var element = getRelevantTableElement(week, person, day)

		// TODO: this can definitely be improved
		element.innerHTML += `<div class="note">
								<div class="main-note" onclick="revealNote(this)">
									<div class="note-container" style="display: none;">
										<div class="note-content">${content}</div>
									</div>
								</div>
								<div class="note-close-button" style="display: none;" onclick="closeNote(this)">X</div>
							  </div>`
		

	} else if (command == "deletenote") {
		var arr = message.split("$")	

		var content, week, person, day

		[_, content, week, person, day] = arr

		var element = getRelevantTableElement(week, person, day)

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
		socket_conn.send("toggle$" + week + "$" + person + "$" + day)
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

	console.log(person, day)
	add_note_button.setAttribute("name", "current$" + person + "$" + day)
	delete_note_button.setAttribute("name", "current$" + person + "$" + day)

})


// note handling
function revealNote(element) {
	element.children[0].style.display = "inline-block"

	element.parentNode.children[1].style.display = "inline-block"
}

function closeNote(element) {
	element.parentNode.children[0].children[0].style.display = "none"
	element.style.display = "none"
}


function add_new_note(name) {
	console.log("added new note")

	var new_note = prompt("What should the note say")

	var week, name, day;

	[week, name, day] = name.split("$")

	if (week == "" || name == "" || day == "") {
		custom_context.style.display = "none"
		return
	}

	console.log(week, name, day)

	custom_context.style.display = "none"

	socket_conn.send("addnote$" + new_note + "$" + week + "$" + name + "$" + day)
}

function edit_note(name) {
	// remove note then add note
}

function delete_note(name) {
	console.log("deleting note...")

	var week, person, day;
	[week, person, day] = name.split("$")

	if (week == "" || person == "" || day == "") {
		custom_context.style.display = "none"
		console.log("remove note:", week, person, day)
		console.log("invalid, aborting")
		return
	}

	console.log("remove note:", week, person, day)

	custom_context.style.display = "none"

	var new_note = getRelevantTableElement(week, person, day).innerHTML

	socket_conn.send("deletenote$" + new_note + "$" + week + "$" + person + "$" + day)
}

// admin panel
function toggle_admin_panel(element) {
	var display = element.parentNode.children[1].style.display

	element.parentNode.children[1].style.display = display == "none" ? "flex" : "none"

	console.log("toggled admin panel")
}

// utility
function getRelevantTableElement(week, person, day) {
	var week_table = document.getElementsByClassName(week + "-week")[0]
	var element = week_table.getElementsByClassName(day)[0].getElementsByClassName(person)[0]

	return element
}


