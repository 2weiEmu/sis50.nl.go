console.log("Loaded index.js")

// NOTE: weekday section
//       ---------------
const weekdayList = ["zo", "ma", "di", "wo", "do", "vr", "za"]
const date = new Date()

setInterval(setWeekday(date), 600000);
function setWeekday(date) {
	var day = weekdayList[date.getDay()]
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

var dayWebsocket = new Websocket("localhost:8000/day-ws")

dayWebsocket.onopen = (event) => {
	// TODO:
}

var gridElList = Array.from(document.getElementsByClassName("rick"))
gridElList = gridElList.concat(Array.from(document.getElementsByClassName("youri")))
gridElList = gridElList.concat(Array.from(document.getElementsByClassName("robert")))
gridElList = gridElList.concat(Array.from(document.getElementsByClassName("milan")))
console.log(gridElList)

for (var i = 0; i < gridElList.length; i++) {
	gridElList[i].addEventListener("click", function (event) {
		var day, person, state

		day = this.classList

		dayWebsocket.send(JSON.stringify({
			"day": day,
			"person": person,
			"state": state,
		}))
	})
}

// NOTE: helper functions
// 	     ----------------

function findDayFromList(classList, dayList) {

}

function findPersonFromDayList
