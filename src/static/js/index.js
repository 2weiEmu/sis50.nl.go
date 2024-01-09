console.log("Loaded index.js")

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
