
// NOTE: shopping websocket section
// 	     ---------------------
// id int
// content string
// action string

var shoppingList = document.getElementById("shop-list")

var shopWebSocket

var shopPonged = true

function connectShop() {
	console.log("connectShop")
	shopWebSocket = new WebSocket(`${secure}://${WS_BASE}/shopWS`, "echo-protocol")
	shopWebSocket.onopen = () => {
		clearInterval(this.timerId)
		console.log("[INFO] Opening Shop Socket")
		shoppingList.innerHTML = ""
		shopWebSocket.send(JSON.stringify({
			"id": "0",
			"content": "",
			"action": "open-shopping"
		}))
		shopWebSocket.onclose = (event) => {
			this.timerId = setInterval(() => {
				connectShop()
			}, 200);
		}
	}

	shopWebSocket.onmessage = async (event) => {
		console.log("[INFO] Received Shop Message")
		var message = JSON.parse(event.data)
		console.log(message)

		if (message.action == "add") {
			var shoppingItem = document.createElement("div")
			shoppingItem.classList.add("shopping-item")
			shoppingItem.id = message.id

			shoppingItem.draggable = true
			shoppingItem.addEventListener("dragstart", handleDragStart)
			shoppingItem.addEventListener("dragend", handleDragEnd)

			shoppingItem.addEventListener("dragenter", handleDragEnter)
			shoppingItem.addEventListener("dragleave", handleDragLeave)
			shoppingItem.addEventListener("dragover", function(ev){ ev.preventDefault() })

			shoppingItem.addEventListener("drop", handleDropOn)

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
		} else if (message.action == "pong") {
			console.log("received PONG in the shop")
			shopPonged = true
		} else if (message.action == "edit") {
			for (var i = 0; i < shoppingList.children.length; i++) {
				if (shoppingList.children[i].id == message.id) {
					shoppingList.children[i].children[0].innerText = message.content
				}
			}
		} else if (message.action == "rearrange") {
			insertInShoppingListAtIndex(message.id, message.content)
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
}

connectShop()

function addItem() {
	var content = document.getElementById("item-name-add").value
	shopWebSocket.send(JSON.stringify({
		"id": "-1",
		"content": content,
		"action": "add"
	}))
}

function editItem(_) {
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

function removeItem(_) {
	var id = this.parentElement.id
	console.log("[INFO] ID to remove:", id)
	shopWebSocket.send(JSON.stringify({
		"id": `${id}`, 
		"content": "",
		"action": "remove"
	}))
}

let dragged;

function handleDragStart(event) {
	console.log("Started dragging")
	dragged = event.target;
}

function handleDragEnd() {
	console.log("Stopped dragging")

	for (let i = 0; i < shoppingList.children.length; i++) {
		shoppingList.children[i].classList.remove("hovered-over")
	}
	dragged = null
}

function handleDragEnter(event) {
	let el = event.target
	if (indexInShoppingList(el) > indexInShoppingList(dragged)) {
		el.classList.add("hovered-over-low")
	}
	else {
		el.classList.add("hovered-over-high")
	}
}

function handleDragLeave(event) {
	let el = event.target
	el.classList.remove("hovered-over-high")
	el.classList.remove("hovered-over-low")
}

function handleDropOn(event) {
	console.log("dropping")
	let el = event.target

	if (dragged === null) return

	console.log(`INDEX IN THE SHOPPING LIST WHERE IT WAS DROPPED MATE: ${indexInShoppingList(el)}`)
	shopWebSocket.send(JSON.stringify({
		"id": dragged.id,
		"content": `${indexInShoppingList(el)}`,
		"action": "rearrange"
	}))
}

function indexInShoppingList(element) {
	let len = shoppingList.children.length
	for (let i = 0; i < len; i++) {
		if (shoppingList.children[i] == element) {
			return i;
		}
	}

	return -1;
}

function insertInShoppingListAtIndex(id, newIndex) {
	let to_remove = document.getElementById(id)
	let target = newShoppingItem(to_remove.id, to_remove.children[0].innerText)
	let oldIndex = indexInShoppingList(to_remove)
	
	if (newIndex >= shoppingList.children.length) {
		shoppingList.appendChild(target)
	} else if (oldIndex > newIndex) {
		shoppingList.insertBefore(target, shoppingList.children[newIndex])
	} else {
		shoppingList.insertBefore(target, shoppingList.children[newIndex].nextSibling)
	}
	to_remove.remove()
}

function newShoppingItem(id, text) {
	var shoppingItem = document.createElement("div")
	shoppingItem.classList.add("shopping-item")
	shoppingItem.id = id

	shoppingItem.draggable = true

	shoppingItem.addEventListener("dragstart", handleDragStart)
	shoppingItem.addEventListener("dragend", handleDragEnd)

	shoppingItem.addEventListener("dragenter", handleDragEnter)
	shoppingItem.addEventListener("dragleave", handleDragLeave)
	shoppingItem.addEventListener("dragover", function(ev){ ev.preventDefault() })

	shoppingItem.addEventListener("drop", handleDropOn)

	var item_desc = document.createElement("p")
	item_desc.innerText = text

	var edit_button = document.createElement("button")
	edit_button.addEventListener("click", editItem)
	edit_button.innerText = "Edit"

	var remove_button = document.createElement("button")
	remove_button.addEventListener("click", removeItem)
	remove_button.innerText = "Remove"


	shoppingItem.appendChild(item_desc)
	shoppingItem.appendChild(edit_button)
	shoppingItem.appendChild(remove_button)

	return shoppingItem
}


function cleanInput() {
	document.getElementById("item-name-add").value = ""
	return false;
}

document.getElementById("logoutButton").addEventListener("click", async () => {
	console.log("logging out")
	const response = await fetch("/logout", {
			method: "POST",
			headers: {
				'Content-Type': 'application/json'
			},
		}
	)

	if (response.status == 200) {
		window.location.href = "/"
	}
})

var notifNote = document.getElementsByClassName("message-notif")[0].children[1]
notifNote.innerHTML = formatMessage(notifNote.innerText)
console.log("Formatted.")

function setDay(el) {
	var day, person, state
	day = clickedOnDay
	person = clickedOnPerson
	state = el.getAttribute("state")

	console.log("Sending message...")
	console.log(day, person, state)
	dayWebsocket.send(JSON.stringify({
		"day": day,
		"person": person,
		"state": state,
	}))
	closeDialog()
}

const SHOPPING_PING_DELAY = 20000
setInterval(() => {
	console.log("new set interval: shopPonged: " + shopPonged)
	if (shopPonged == false) {
		console.log("failed to receive a PONG for Shop, reloading!")
		location.reload()
	} else {
		shopPonged = false
		console.log("PINGING FOR SHOP")
		shopWebSocket.send(JSON.stringify({
			"id": "0",
			"content": "something",
			"action": "ping"
		}))
	}
}, SHOPPING_PING_DELAY)
