import { 
	getRelevantTableElement, current_week_table, next_week_table, people, days,
	send_day_toggle

} from "./koken-helper.js" 

// variables

var socket_conn = new WebSocket("ws://localhost:8000/koken-ws")

var custom_context = document.getElementById("custom-context-menu")
var add_note_button = document.getElementById("add-note-button")
var delete_note_button = document.getElementById("remove-note-button")
var berichte_list = document.getElementById("berichte-list")

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

		var element = getRelevantTableElement(message.week, message.person, message.day)
		element.innerHTML = message.currentState + element.innerHTML.slice(1)


	} else if (command == "post-bericht") {
		const para = document.createElement("p")
		para.appendChild(document.createTextNode(message.currentState))
		berichte_list.prepend(para)

		// if there are more than 5 berichte, we must remove the most bottom one
		// this is just part of the functionality
		if (berichte_list.children.length > 5) {
			socket_conn.send(JSON.stringify({
				command: "del-bericht"
			}))
		}
		

	} else if (command == "del-bericht") {
		console.log("removing bericht...")
		var remove_child = berichte_list.childNodes[berichte_list.children.length - 1]
		console.log(remove_child)
		berichte_list.removeChild(remove_child)


	} else if (command == "addnote") {
		var elem = getRelevantTableElement(message.week, message.person, message.day)	
		add_note(elem, message.optid, message.currentState)


	} else if (command == "deletenote") {

	}

}


function add_note(element, id, currentState) {
	var note_div = document.createElement('div')
	note_div.style.backgroundColor = "yellow"
	note_div.innerText = currentState
	note_div.id = `n${id}`

	var note_button = document.createElement('button')
	note_button.innerText = "X"
	note_button.id = `${id}`

	note_button.addEventListener("click", () => { remove_note(id) })

	element.appendChild(note_div)
	element.appendChild(note_button)
}

function remove_note(id) {
	var button = document.getElementById(id)
	var div = document.getElementById(`n${id}`)

	button.parentNode.removeChild(button)
	div.parentNode.removeChild(div)
}

// day toggling mechanisms
current_week_table.addEventListener("mousedown", function(ev) {
	send_day_toggle(window, socket_conn, ev, "current")
})

next_week_table.addEventListener("mousedown", function(ev) {
	send_day_toggle(window, socket_conn, ev, "next")
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
