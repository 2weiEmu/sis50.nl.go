console.log("Loaded calendarHandler.js")

const stateList = ["present", "absent", "cooking", "uncertain", "maybe-cooking", "cant-cook"]
const altTextList = ["Present", "Absent", "Cooking", "Uncertain if Present", "Maybe Cooking", "Can't Cook"]

var gridElList = Array.from(document.getElementsByClassName("rick"))
gridElList = gridElList.concat(Array.from(document.getElementsByClassName("youri")))
gridElList = gridElList.concat(Array.from(document.getElementsByClassName("robert")))
gridElList = gridElList.concat(Array.from(document.getElementsByClassName("milan")))
console.log(gridElList)

let clickedOnDay
let clickedOnPerson

var dayWebsocket = new WebSocket(`${secure}://${WS_BASE}/dayWS`, "echo-protocol")

function makeCalMsg(day, person, state) {
	return JSON.stringify({
		"day": day,
		"person": person,
		"state": state
	})
}

dayWebsocket.onopen = () => {
	dayWebsocket.send(makeCalMsg("", "", "open-calendar"))
}

function handleGridClick() {
	let day = this.getAttribute("data-day")
	let person = this.getAttribute("data-person")
	let state = this.childNodes[0].getAttribute("data-state")
	dayWebsocket.send(makeCalMsg(day, person, state))
}

function handleGridMenu(event) {
	if (dialogOpen) {
		closeDialog()
	}
	else {
		clickedOnDay = this.getAttribute("data-day")
		clickedOnPerson = this.getAttribute("data-person")
		openDialogAt(event.clientX, event.clientY)
	}
	event.preventDefault()
}

for (var i = 0; i < gridElList.length; i++) {
	gridElList[i].addEventListener("click", handleGridClick)
	gridElList[i].addEventListener("contextmenu", handleGridMenu)
	gridElList[i].childNodes[0].setAttribute("draggable", false)
}

dayWebsocket.onmessage = async function(event) {
	var message = JSON.parse(event.data)

	if (message.state != "open-calendar") {
		// Update the day's state
		var el = document.getElementsByClassName(`${message.person} ${message.day}`)[0]
		let state_image = el.childNodes[0]

		// TODO: move this to CSS
		state_image.style.width = "0"
		state_image.style.height = "0"
		state_image.style.marginRight = "50%"
		state_image.style.marginLeft = "50%"
		state_image.style.marginTop = "30%"
		state_image.style.marginBottom = "50%"
		await sleep(300)

		state_image.setAttribute("data-state", message.state) 
		state_image.src = "/images/" + message.state + ".svg"

		state_image.style.width = "100%"
		state_image.style.marginRight = "0"
		state_image.style.marginLeft = "0"
		state_image.style.marginTop = "8px"
		state_image.style.marginBottom = "0"
		state_image.style.height = "80%"
		await sleep(300)

		console.log(`[INFO] message.state: ${message.state}`)
		var i = stateList.findIndex((item) => { return item == message.state })
		state_image.title = altTextList[i]
		
	} else {
		var days = message.day.split("/")
		for (var i = 0; i < days.length; i++) {
			day_states = days[i].split("")
			var day = weekdayList[((i + 1) % 7)]
			console.log("Open day:", day)
			for (var j = 0; j < day_states.length; j++) {
				var person = personList[j]
				var newState = stateList[day_states[j]]
				var el = document.getElementsByClassName(`${person} ${day}`)[0]
				console.log(el)
				console.log(`[INFO] day_states[j] ${day_states[j]}`)
				el.childNodes[0].setAttribute("data-state", newState) 
				el.childNodes[0].src = "/images/" + newState + ".svg"
				el.childNodes[0].title = altTextList[day_states[j]]
			}
		}
	}
}

