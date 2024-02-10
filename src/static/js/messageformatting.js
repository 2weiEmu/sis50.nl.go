var symbolMap = {
	"*": "b",
	"_": "i",
	"~": "s",
}

// TODO: reformat this please
function formatMessage(text) {
	var result = ""
	var boldOpen = false
	var italicOpen = false
	var strikeOpen = false
	text = text + " "

	for (var i = 0; i < text.length - 1; i++) {

		if (text[i] == "*" && text[i + 1] != " " && !boldOpen) {
			boldOpen = true
			result += "<b>"
		}
		else if (text[i] == "*" && boldOpen) {
			boldOpen = false
			result += "</b>"
		}
		else if (text[i] == "_" && text[i + 1] != " " && !italicOpen) {
			italicOpen = true
			result += "<i>"
		}
		else if (text[i] == "_" && italicOpen) {
			italicOpen = false
			result += "</i>"
		}
		else if (text[i] == "~" && text[i + 1] != " " && !strikeOpen) {
			strikeOpen = true
			result += "<s>"
		}
		else if (text[i] == "~" && strikeOpen) {
			strikeOpen = false
			result += "</s>"
		}
		else {
			result += text[i]
		}
	}

	return result
}
