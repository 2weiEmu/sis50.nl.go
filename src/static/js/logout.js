console.log("logout.js loaded");

document.getElementById("logoutButton").addEventListener("click", async () => {
	console.log("logging out")
	const response = await fetch("/logout", {
			method: "POST",
			headers: {
				'Content-Type': 'application/json'
			},
		}
	)

	if (response.status == 200) {
		window.location.href = "/"
	}
})
