console.log("Loaded messages.js")

document.getElementById("bg-button").addEventListener("click", bgMenu)

if (localStorage.getItem("sis50-background") === null) {
	localStorage.setItem("sis50-background", "")
}

// TODO: this does not actually work yet - has to be fixed before deployment
var ws_args = document.currentScript.getAttribute("args")
var argv = ws_args.split(" ")

var WS_BASE = argv[0]

window.onload = (event) => {
	// Make GET request to receive first page
	$.ajax({
		url: `http://localhost:8000/api/messages/0`,
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
				msgPar.innerText = messagePageZero.messages[i]
				el.appendChild(msgPar)
				msgList.appendChild(el)
			}

			if (messagePageZero.messages.length < 5) {
				loadMoreItems(document.getElementById("load-button"))
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
		url: "http://localhost:8000/api/messages",
		type: "POST",
		dataType: "JSON",
		data: data,
		success: function(data) { alert("yes") },
		error: function(req, error) { alert("no") }
	})
}

function loadMoreItems(button) {
	var pageNumber = button.getAttribute("data-page-number")
	pageNumber++

	$.ajax({
		url: `http://localhost:8000/api/messages/${pageNumber}`,
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
				msgPar.innerText = messagePageZero.messages[i]
				el.appendChild(msgPar)
				msgList.appendChild(el)
			}
			button.setAttribute("data-page-number", pageNumber)
		},
		error: function(req, error) {
			alert("Messages could not be loaded.")
		}
	})

}
