console.log("[INFO] Loaded messages.js")

document.getElementById("msg-form").addEventListener("submit", addMessage)
document.getElementById("bg-button").addEventListener("click", bgMenu)

if (localStorage.getItem("sis50-background") === null) {
	localStorage.setItem("sis50-background", "")
}

var ws_args = document.currentScript.getAttribute("args")
var argv = ws_args.split(" ")

var WS_BASE = argv[0]
console.log(`[INFO] ${WS_BASE}`)

window.onload = (event) => {
	$.ajax({
		url: `http://${WS_BASE}/api/messages/0`,
		type: 'GET',
		dataType: 'json',
		CORS: true,
		headers: {
			'Access-Control-Allow-Origin': '*',
		},
		success: function(data) {
			console.log(`[INFO] Ajax Success to /api/messages/0 with data:\n${data}`)
			var messagePageZero = data

			var msgList = document.getElementById("msg-list")

			for (var i = messagePageZero.messages.length - 1; i > -1; i--) {
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
			alert(`[ERROR] ${error} -> Loading first page of Messages`)
		}
	})
}

function addMessage(event) {
	var msg = document.getElementById("msg-content").value
	var data = JSON.stringify({
		"message": msg
	})

	$.ajax({
		url: `http://${WS_BASE}/api/messages`,
		type: "POST",
		data: data,
		async: false, // ok apparently sync here is bad (it's deprecated, but it makes this work on firefox which I like)
		success: function(data) { 
			return false
		},
		error: function(req, error) { 
			alert(`[ERROR] ${error} -> when trying to add a message`)
			return false
		}
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
			console.log(`[INFO] Loading More items on PageNumber: ${pageNumber} | data:\n${data}`)
			var messagePageZero = data

			var msgList = document.getElementById("msg-list")

			for (var i = messagePageZero.messages.length - 1; i > -1; i--) {
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

