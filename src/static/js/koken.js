
var current_week_table = document.getElementsByClassName("current-week")[0];
var next_week_table = document.getElementsByClassName("next-week")[0];

var socket_conn = new WebSocket("ws://localhost:8000/koken-ws")

socket_conn.onmessage = function(event) {
	console.log(event.data)
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
