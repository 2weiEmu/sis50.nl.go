var symbols = [
	["*", "b", false],
	["_", "i", false],
	["~", "s", false],
]

// TODO: reformat this please
function formatMessage(text) {
	var result = ""
	var flag = false
	text = text + " "

	for (var i = 0; i < text.length - 1; i++) {

		for (var j = 0; j < symbols.length; j++) {
			var targetSymbol = symbols[j][0]
			var symbolOpen = symbols[j][2]
			var insertSymbol = symbols[j][1]

			if (text[i] == targetSymbol && text[i + 1] != " " && !symbolOpen) {
				symbols[j][2] = true
				result += `<${insertSymbol}>`
				flag = true
			}
			else if (text[i] == targetSymbol && symbolOpen) {
				symbols[j][2] = false
				result += `</${insertSymbol}>`
				flag = true
			}
		}

		if (!flag) {
			result += text[i]
		}
		flag = false
	}
	return result
}
