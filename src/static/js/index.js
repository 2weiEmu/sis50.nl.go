console.log("Loaded index.js")

function sleep(ms = 0) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

document.getElementById("bg-button").addEventListener("click", bgMenu)

if (localStorage.getItem("sis50-background") === null) {
	localStorage.setItem("sis50-background", "")
}

// NOTE: getting arguments, setting state
// --------------------------------

var ws_args = document.currentScript.getAttribute("args")
var argv = ws_args.split(" ")

var WS_BASE = argv[0]
var security = argv[1]

var secure = "ws"

if (security != "none") {
	secure = "wss"
}

console.log(`secure: ${secure}`)
console.log(`WS_BASE: ${WS_BASE}`)

// NOTE: weekday section
//       ---------------
const weekdayList = ["zo", "ma", "di", "wo", "do", "vr", "za"]
const personList = ["rick", "youri", "robert", "milan"]
const date = new Date()

let day = weekdayList[date.getDay()]
let yesterday = weekdayList[(date.getDay() + 6) % weekdayList.length]

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

// NOTE: day websocket section
// 	     ---------------------
// to communicate use:
// day: string 
// person: string 
// state: string

var stateList = ["present", "absent", "cooking", "uncertain", "maybe-cooking", "cant-cook"]
var altTextList = ["Present", "Absent", "Cooking", "Uncertain if Present", "Maybe Cooking", "Can't Cook"]

var dayWebsocket = new WebSocket(`${secure}://${WS_BASE}/dayWS`, "echo-protocol")

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

dayWebsocket.onmessage = async function(event) {
	console.log("Received Message")
	var message = JSON.parse(event.data)
	console.log(message)

	if (message.state != "open-calendar") {
		// Update the day's state
		var el = document.getElementsByClassName(`${message.person} ${message.day}`)[0]
		let state_image = el.childNodes[0]

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


// NOTE: shopping websocket section
// 	     ---------------------
// id int
// content string
// action string

var shoppingList = document.getElementById("shop-list")

var shopWebSocket = new WebSocket(`${secure}://${WS_BASE}/shopWS`, "echo-protocol")

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

shopWebSocket.onmessage = async (event) => {
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
				var h = shoppingList.children[i].clientHeight
				shoppingList.children[i].style.height = h + "px";
				await new Promise(r => setTimeout(r, 10));
				for (var j = 0; j < shoppingList.children[i].children.length; j++) {
					shoppingList.children[i].children[j].innerText = ""
				}
				shoppingList.children[i].style.height = "1px";
				await new Promise(r => setTimeout(r, 200));
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

var notifNote = document.getElementsByClassName("message-notif")[0].children[1]
notifNote.innerHTML = formatMessage(notifNote.innerText)
console.log("Formatted.")
