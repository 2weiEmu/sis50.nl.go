
var current_week_table = document.getElementsByClassName("current-week")[0];
var next_week_table = document.getElementsByClassName("next-week")[0];

var socket_conn = new WebSocket("ws://localhost:8000/koken-ws")

socket_conn.onmessage = function(event) {

	var message = event.data

	var arr = message.split("$")

	var state, week, person, day;
	
	[state, week, person, day] = arr

	console.log(state, week, person, day)

	if (state == "E") { state = " " }

	var week_table = document.getElementsByClassName(week + "-week")[0]



	var element = week_table.getElementsByClassName(day)[0].getElementsByClassName(person)[0]

	element.innerHTML = state;
	
}

current_week_table.addEventListener("click", function(ev) {

	var target = ev.target
	var person = ev.target.className
	var day = ev.target.parentNode.className

	if (person == "" || day == "") {
		return
	} else {
		console.log(person, day)
		socket_conn.send("toggle$current$" + person + "$" + day)
	}
})

next_week_table.addEventListener("click", function(ev) {

	var target = ev.target
	var person = ev.target.className
	var day = ev.target.parentNode.className

	if (person == "" || day == "") {
		return
	} else {
		console.log(person, day)
		socket_conn.send("toggle$next$" + person + "$" + day)
	}
})
