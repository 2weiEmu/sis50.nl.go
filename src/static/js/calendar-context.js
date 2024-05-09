console.log("Loaded calendar-context.js")

let dialog = document.getElementById("right-click-menu")

function openDialogAt(x, y) {
	dialog.style.display = "block"
	dialog.style.top = y + "px"
	dialog.style.left = x + "px"
}

function closeDialog() {
	dialog.style.display = "none"
}

function setDay(el) {
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
}

// var stateList = ["present", "absent", "cooking", "uncertain", "maybe-cooking", "cant-cook"]

