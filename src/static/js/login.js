console.log("Loaded login.js")

let user_input = document.getElementById("login-user")
let pass_input = document.getElementById("login-password")

async function login() {
	let username = user_input.value
	let password = pass_input.value
	
	const response = await fetch("/login", {
			method: "POST",
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				username: username,
				password: password
			})
		}
	)

	if (response.status == 200) {
		window.location.href = "/"
	}
	else {
		alert(await response.text())
	}
	console.log(response.status)
}
