export function handleToggleMessage(message) {
	var arr = message.split("$")

	var state = arr[1]
	var e, week, person, day;
	
	[e, state, week, person, day] = arr

	console.log(state, week, person, day)

	if (state == "E") { state = "_" }

	var element = getRelevantTableElement(week, person, day)

	// for now we jankily replace the first character, as that is the only thing that should have to change
	element.innerHTML = state + element.innerHTML.slice(1)
}


export function getRelevantTableElement(week, person, day) {
	var week_table = document.getElementsByClassName(week + "-week")[0]
	var element = week_table.getElementsByClassName(day)[0].getElementsByClassName(person)[0]

	return element
}
