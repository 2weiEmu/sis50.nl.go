console.log("Loaded messages.js")

document.getElementById("bg-button").addEventListener("click", bgMenu)

if (localStorage.getItem("sis50-background") === null) {
	localStorage.setItem("sis50-background", "")
}

// TODO: this does not actually work yet - has to be fixed before deployment
var ws_args = document.currentScript.getAttribute("args")
var argv = ws_args.split(" ")

var WS_BASE = argv[0]
console.log(WS_BASE)

window.onload = (event) => {
	// Make GET request to receive first page
	$.ajax({
		url: `http://${WS_BASE}/api/messages/0`,
		type: 'GET',
		dataType: 'json',
		CORS: true,
		headers: {
			'Access-Control-Allow-Origin': '*',
		},
		success: function(data) {
			console.log(data)
			var messagePageZero = data

			var msgList = document.getElementById("msg-list")

			for (var i = messagePageZero.messages.length - 1; i > -1; i--) {
				console.log(messagePageZero.messages[i])
				var el = document.createElement("div")
				el.classList.add("msg")
				var msgPar = document.createElement("p")
				msgPar.innerHTML = formatMessage(messagePageZero.messages[i]);
				el.appendChild(msgPar)
				msgList.appendChild(el)
			}

			if (messagePageZero.messages.length < 5) {
				loadMoreItems(document.getElementById("load-button"), true)
			}
		},
		error: function(req, error) {
			alert("Messages could not be loaded.")
		}
	})
}

function addMessage() {
	var msg = document.getElementById("msg-content").value
	var data = JSON.stringify({
		"message": msg
	})

	$.ajax({
		url: `http://${WS_BASE}/api/messages`,
		type: "POST",
		data: data,
		success: function(data) {},
		error: function(req, error) { alert(error); return false; }
	})
}

function loadMoreItems(button, flag = false) {
	var pageNumber = button.getAttribute("data-page-number")
	pageNumber++

	$.ajax({
		url: `http://${WS_BASE}/api/messages/${pageNumber}`,
		type: 'GET',
		dataType: 'json',
		CORS: true,
		headers: {
			'Access-Control-Allow-Origin': '*',
		},
		success: function(data) {
			console.log(data)
			var messagePageZero = data

			var msgList = document.getElementById("msg-list")

			for (var i = messagePageZero.messages.length - 1; i > -1; i--) {
				console.log(messagePageZero.messages[i])
				var el = document.createElement("div")
				el.classList.add("msg")
				var msgPar = document.createElement("p")
				msgPar.innerHTML = formatMessage(messagePageZero.messages[i]);
				el.appendChild(msgPar)
				msgList.appendChild(el)
			}
			button.setAttribute("data-page-number", pageNumber)
		},
		error: function(req, error) {
			if (!flag && error != "parsererror") {
				alert("Messages could not be loaded.")
			} 

		}
	})

}

var symbolMap = {
	"*": "b",
	"_": "i",
	"~": "s",
}

function formatMessage(text) {
	var open = -1;
	var close = -1;

	var closeOpenList = [];
	var len = text.length;
	text += " " // add one extra, for the sake of convenience. Not performant.

	for (var i = 0; i < len; i++) {
		if (text[i] == "*" && text[i+1] != " " && open == -1) {
			open = i;
		}
		else if (open != -1 && text[i] == "*") {
			close = i;

			closeOpenList.push(open);
			closeOpenList.push(close);
		}
	}

	console.log("closeOpenList: ", closeOpenList)
	console.log("text: ", text)

	var pieces = [];
	var prev = 0;
	for (var i = 0; i < closeOpenList.length; i++) {
		var pos = closeOpenList[i];
		var slice = text.slice(prev, pos)
		console.log("slice: ", slice, "at i: ", i)
		pieces.push(slice)
		prev = pos + 1;
	}

	pieces.push(text.slice(prev, text.length - 1))

	console.log("pieces: ", pieces)

	var openTag = true;
	var result = "";
	for (var i = 0; i < pieces.length - 1; i++) {
		if (openTag) {
			result += pieces[i] + "<b>";
		} else {
			result += pieces[i] + "</b>"
		}
		openTag = !openTag
	}

	if (closeOpenList.length != 0) {
		result += pieces[pieces.length - 1]
	} else {
		result = text;
	}

	console.log("result: ", result)
	
	return result
}
