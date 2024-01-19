console.log("Loaded messages.js")

document.getElementById("bg-button").addEventListener("click", bgMenu)

if (localStorage.getItem("sis50-background") === null) {
	localStorage.setItem("sis50-background", "")
}

// TODO: this does not actually work yet - has to be fixed before deployment
var ws_args = document.currentScript.getAttribute("args")
var argv = ws_args.split(" ")

var WS_BASE = argv[0]
