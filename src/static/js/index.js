console.log("Loaded index.js")

document.getElementById("bg-button").addEventListener("click", bgMenu)

if (localStorage.getItem("sis50-background") === null) {
	localStorage.setItem("sis50-background", "")
}

// NOTE: getting arguments, setting state
//		 --------------------------------

var ws_args = document.currentScript.getAttribute("args")
var argv = ws_args.split(" ")

var WS_BASE = argv[0]
console.log(`WS_BASE: ${WS_BASE}`)

// NOTE: weekday section
//       ---------------
const weekdayList = ["zo", "ma", "di", "wo", "do", "vr", "za"]
const personList = ["rick", "youri", "robert", "milan"]
const date = new Date()

setInterval(setWeekday(date), 600000);
function setWeekday(date) {
	var day = weekdayList[date.getDay()]
	var old_day = weekdayList[Math.abs(date.getDay() - 6)]
	var old_el = document.getElementsByClassName(`day ${old_day}`)[0]
	old_el.classList.remove("selecteled")

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
var altTextList = ["Present", "Absent", "Cooking", "Uncertain if Present", "Maybe Cooking", "Can't Cook"]

console.log(`ws://${WS_BASE}/dayWS`)
var dayWebsocket = new WebSocket(`ws://${WS_BASE}/dayWS`, "echo-protocol")

dayWebsocket.onopen = (event) => {
	dayWebsocket.send(JSON.stringify({
		"day": "",
		"person": "",
		"state": "open-calendar"
	}))
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

	gridElList[i].childNodes[0].setAttribute("draggable", false)
}

dayWebsocket.onmessage = function(event) {
	console.log("Received Message")
	var message = JSON.parse(event.data)
	console.log(message)

	if (message.state != "open-calendar") {
		var el = document.getElementsByClassName(`${message.person} ${message.day}`)[0]
		console.log(el)
		console.log(`${message.person} ${message.day}`)
		el.childNodes[0].setAttribute("data-state", message.state) 
		el.childNodes[0].src = "/images/" + message.state + ".svg"
		var i = stateList.findIndex((item) => { item == message.state })
		console.log(i)
		el.childNodes[0].title = altTextList[i]
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
				el.childNodes[0].setAttribute("data-state", newState) 
				el.childNodes[0].src = "/images/" + newState + ".svg"
				el.childNodes[0].title = altTextList[day_states[j]]
			}
		}
	}
}


// NOTE: shopping websocket section
// 	     ---------------------
// id int
// content string
// action string

var shoppingList = document.getElementById("shop-list")

var shopWebSocket = new WebSocket(`ws://${WS_BASE}/shopWS`, "echo-protocol")

function addItem() {
	var content = document.getElementById("item-name-add").value
	shopWebSocket.send(JSON.stringify({
		"id": "-1",
		"content": content,
		"action": "add"
	}))
}

function editItem(event) {
	var old_content = this.parentElement.children[0].innerText
	var new_content = window.prompt("Edit the name", `${old_content}`)
	var id = this.parentElement.id
	console.log("[INFO] ID to remove:", id)
	shopWebSocket.send(JSON.stringify({
		"id": `${id}`, 
		"content": new_content,
		"action": "edit"
	}))
}

function removeItem(event) {
	var id = this.parentElement.id
	console.log("[INFO] ID to remove:", id)
	shopWebSocket.send(JSON.stringify({
		"id": `${id}`, 
		"content": "",
		"action": "remove"
	}))
}

shopWebSocket.onopen = () => {
	console.log("[INFO] Opening Shop Socket")
	shopWebSocket.send(JSON.stringify({
		"id": "0",
		"content": "",
		"action": "open-shopping"
	}))
}

shopWebSocket.onmessage = (event) => {
	console.log("[INFO] Received Shop Message")
	var message = JSON.parse(event.data)
	console.log(message)

	if (message.action == "add") {
		var shoppingItem = document.createElement("div")
		shoppingItem.classList.add("shopping-item")
		shoppingItem.id = message.id

		var item_desc = document.createElement("p")
		item_desc.innerText = message.content

		var edit_button = document.createElement("button")
		edit_button.addEventListener("click", editItem)
		edit_button.innerText = "Edit"

		var remove_button = document.createElement("button")
		remove_button.addEventListener("click", removeItem)
		remove_button.innerText = "Remove"


		shoppingItem.appendChild(item_desc)
		shoppingItem.appendChild(edit_button)
		shoppingItem.appendChild(remove_button)
		shoppingList.appendChild(shoppingItem)

	} else if (message.action == "edit") {
		for (var i = 0; i < shoppingList.children.length; i++) {
			if (shoppingList.children[i].id == message.id) {
				shoppingList.children[i].children[0].innerText = message.content
			}
		}
	} else {
		// remove the item by id
		for (var i = 0; i < shoppingList.children.length; i++) {
			if (shoppingList.children[i].id == message.id) {
				shoppingList.children[i].remove()
				break
			}
		}
	}
}

function cleanInput() {
	document.getElementById("item-name-add").value = ""
	return false;
}

// Background images

