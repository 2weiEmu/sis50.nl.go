console.log("Loaded calendar-context.js")

let dialog = document.getElementById("right-click-menu")

let dialogOpen = false

function openDialogAt(x, y) {
	dialog.style.display = "block"
	dialog.style.top = y + "px"
	dialog.style.left = x + "px"
	dialogOpen = true
}

function closeDialog() {
	dialog.style.display = "none"
	dialogOpen = false
}


// var stateList = ["present", "absent", "cooking", "uncertain", "maybe-cooking", "cant-cook"]

