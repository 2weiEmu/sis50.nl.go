console.log("Loaded index.js")

// NOTE: getting arguments, setting state
//		 --------------------------------

var arguments = document.currentScript.getAttribute("args")
var argv = arguments.split(" ")

var WS_BASE = argv[0]
console.log(`WS_BASE: ${WS_BASE}`)

// NOTE: weekday section
//       ---------------
const weekdayList = ["zo", "ma", "di", "wo", "do", "vr", "za"]
const date = new Date()

setInterval(setWeekday(date), 600000);
function setWeekday(date) {
	var day = weekdayList[date.getDay()]
	console.log(`Got day: ${day}`)
	var el = document.getElementsByClassName(`day ${day}`)[0]
	el.classList.add("selected")
	console.log(`Set Weekday`)
}

// NOTE: day websocket section
// 	     ---------------------
// to communicate use:
// day: string 
// person: string 
// state: string

var stateList = ["present", "absent", "cooking", "uncertain", "maybe-cooking", "cant-cook"]

console.log(`${WS_BASE}/dayWS`)
var dayWebsocket = new WebSocket(`${WS_BASE}/dayWS`, "echo-protocol")

dayWebsocket.onopen = (event) => {
	// TODO:
}

var gridElList = Array.from(document.getElementsByClassName("rick"))
gridElList = gridElList.concat(Array.from(document.getElementsByClassName("youri")))
gridElList = gridElList.concat(Array.from(document.getElementsByClassName("robert")))
gridElList = gridElList.concat(Array.from(document.getElementsByClassName("milan")))
console.log(gridElList)

for (var i = 0; i < gridElList.length; i++) {
	gridElList[i].addEventListener("click", function (event) {
		var day, person, state

		day = this.getAttribute("data-day")
		person = this.getAttribute("data-person")
		state = this.childNodes[0].getAttribute("data-state")

		console.log("Sending message...")
		console.log(day, person, state)
		dayWebsocket.send(JSON.stringify({
			"day": day,
			"person": person,
			"state": state,
		}))
	})
}

dayWebsocket.onmessage = function(event) {
	console.log("Received Message")
	var message = JSON.parse(event.data)
	console.log(message)

	var el = document.getElementsByClassName(`${message.person} ${message.day}`)[0]
	console.log(el)
	el.childNodes[0].setAttribute("data-state", message.state) 
	el.childNodes[0].src = "/images/" + message.state + ".svg"
}
