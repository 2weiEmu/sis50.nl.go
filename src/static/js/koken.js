
var current_week_table = document.getElementsByClassName("current-week")[0];
var next_week_table = document.getElementsByClassName("next-week")[0];

var socket_conn = new WebSocket("ws://localhost:8000/koken-ws")

socket_conn.onopen = function(event) {
	socket_conn.send("open$")
}

socket_conn.onmessage = function(event) {

	var message = event.data
	var arr = message.split("$")

	var state = arr[0]
	var week, person, day;
	
	[state, week, person, day] = arr

	console.log(state, week, person, day)

	if (state == "E") { state = " " }

	var week_table = document.getElementsByClassName(week + "-week")[0]
	var element = week_table.getElementsByClassName(day)[0].getElementsByClassName(person)[0]

	element.innerHTML = state;

	
}

current_week_table.addEventListener("mousedown", function(ev) {

	var target = ev.target
	var person = ev.target.className
	var day = ev.target.parentNode.className

	console.log(ev.button)

	if (ev.button != 0) { return }

	if (person == "" || day == "") {
		return
	} else {
		console.log(person, day)
		socket_conn.send("toggle$current$" + person + "$" + day)
	}

})

document.addEventListener("keyup", function(ev) {
	if (ev.key == "Escape") {
		console.log("escape key pressed")
		context_menu.style.display = "none"
	}
})

var context_menu = document.getElementById("custom-context-menu")

current_week_table.addEventListener("contextmenu", function(ev) {
	ev.preventDefault()

	console.log(ev.clientX);
	console.log(ev.clientY);

	context_menu.style.top = ev.clientY + "px"
	context_menu.style.left = ev.clientX + "px"
	context_menu.style.display = "inline-block";

})


function add_new_note() {
	console.log("added new note")

	var new_note = prompt("What should the note say")
}


next_week_table.addEventListener("mousedown", function(ev) {

	var target = ev.target
	var person = ev.target.className
	var day = ev.target.parentNode.className

	if (ev.button != 0) { return }

	if (person == "" || day == "") {
		return
	} else {
		console.log(person, day)
		socket_conn.send("toggle$next$" + person + "$" + day)
	}
})
