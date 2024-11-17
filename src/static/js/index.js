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
const personList = ["rick", "youri", "robert", "milan"]

let day = weekdayList[date.getDay()]
let yesterday = weekdayList[(date.getDay() + 6) % weekdayList.length]

