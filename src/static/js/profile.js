console.log("loaded profile.js")

const profile_picture_form = document.getElementById("profile-pic-form")

async function uploadProfilePicture(event) {
	event.preventDefault()
	const formData = new FormData(profile_picture_form)

	const response = await fetch("/profile", {
		method: "POST",
		body: formData
	})

	if (response.status != 200) {
		alert("failed to upload image:", response.text())
	}
	else {
		window.location.href = window.location.href
	}
}
