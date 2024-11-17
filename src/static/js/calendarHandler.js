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

const weekdayList = ["zo", "ma", "di", "wo", "do", "vr", "za"]
const date = new Date()
setWeekday(date)
// 600_000 = 10 minutes
setInterval(setWeekday(date), 600_000);
function setWeekday(date) {
	let days = document.getElementsByClassName("day")
	for (let i = 0; i < 7; i++) {
		let set = weekdayList[(date.getDay() + 4 + i) % weekdayList.length]
		set = set[0].toUpperCase() + set[1] + "."
		days[i].innerHTML = `<p>${set}</p>`
	}
}


var dayWebsocket

function connect() {
	console.log("connect")
	dayWebsocket = new WebSocket(`${secure}://${WS_BASE}/dayWS`, "echo-protocol")

	dayWebsocket.iFrame = () => {
		// fetch all of the new grid elements and display them
		dayWebsocket.send(makeCalMsg("", "", "OPEN"))
	}

	dayWebsocket.onmessage = async function(event) {
		const message = JSON.parse(event.data)
		
		if (message.state == "OPEN") {
			// here the serial will be in "person"
			console.log(`Received an OPEN state with the following serial: ${message.person}`)

			const serialList = message.person.split("|")
			serialList.pop()

			for (let i = 0; i < gridElList.length; i++)	{
				const element = gridElList[i]
				const state_image = element.children[0]

				const current_day_index = Math.floor(i / 7)
				const offset_index = i % 7

				const new_state_index = Number(serialList[offset_index][current_day_index])
				const new_state = stateList[new_state_index]

				console.log(`OPEN: element: ${element}, i: ${i}, current_day_index: ${current_day_index}, new_state: ${new_state}`)

				state_image.setAttribute("data-state", new_state)
				state_image.src = "/images/" + new_state + ".svg"

				state_image.title = altTextList[new_state_index]
			}
			

			return
		}

		// Update the day's state
		var el = document.getElementsByClassName(`${message.person} ${message.day}`)[0]
		let state_image = el.children[0]

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

		// some CSS stuff
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
	}


	dayWebsocket.onopen = () => {
		dayWebsocket.onclose = (event) => {
			this.timerId = setInterval(() => {
				connect()
			}, 300);
		}

		clearInterval(this.timerId)
		dayWebsocket.iFrame()
		setWeekday(new Date())

	}

}

connect()


function makeCalMsg(day, person, state) {
	return JSON.stringify({
		"day": day,
		"person": person,
		"state": state
	})
}

function handleGridClick() {
	const day = this.getAttribute("data-day")
	const person = this.getAttribute("data-person")
	const state = this.children[0].getAttribute("data-state")
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
	gridElList[i].children[0].setAttribute("draggable", false)
	console.log(gridElList[i])
}

