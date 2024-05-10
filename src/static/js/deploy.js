console.log("Loaded deploy.js")

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
