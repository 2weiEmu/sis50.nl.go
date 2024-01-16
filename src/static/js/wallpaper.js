function bgMenu(event) {
	console.log("Opening background selection menu...")

	var bg_menu = document.getElementById("bg-menu")
	bg_menu.style.display = "block"

	var bg_confirm = document.getElementById("bg-confirm")
	bg_confirm.setAttribute("data-selected", "")

	console.log("Adding Background Image Listeners...")
	var background_image_list = document.getElementById("bgl")
	console.log(background_image_list)

	document.getElementById("add-new-bg").addEventListener("click", addbg)
	
	updateImgList()

	document.getElementById("bg-confirm").addEventListener("click", confirmBackground)
	document.getElementById("bg-remove").addEventListener("click", removeBackground)


}

function updateImgList() {
	var background_image_list = document.getElementById("bgl")
	var backgrounds = localStorage.getItem("sis50-background").split("^")

	console.log(localStorage.getItem("sis50-background").split("^"))

	background_image_list.innerHTML = ""

	for (var i = 0; i < backgrounds.length; i++) {

		if (backgrounds[i] == "") continue
		var new_img = document.createElement("img")
		new_img.src = backgrounds[i]
		background_image_list.appendChild(new_img)
	}

	for (var i = 0; i < background_image_list.children.length; i++) {
		var background_image = background_image_list.children[i]
		background_image.addEventListener("click", selectBackgroundImage)
	}
}

function selectBackgroundImage(event) {
	console.log("selecting background")
	var bg_confirm = document.getElementById("bg-confirm")
	bg_confirm.setAttribute("data-selected", this.src)
	setBackground(this.src)
}

function setBackground(url) {
	var bod = document.getElementById("mbody")
	console.log(url)
	if (url == "") {
		bod.style.backgroundImage = ""
	} else {
		bod.style.backgroundImage = `url(${url})`
		bod.style.backgroundSize = "100%"
	}
}

function confirmBackground(event) {
	console.log("confirming background")
	setBackground(this.getAttribute("data-selected"))

	var bg_menu = document.getElementById("bg-menu")
	bg_menu.style.display = "none"
	document.getElementById("bg-confirm").removeEventListener("click", confirmBackground)
}

function addbg(event) {
	new_bg = document.getElementById("new-bg").value
	if (new_bg == "") {
		return
	}
	var bg_list = localStorage.getItem("sis50-background").split("^")
	bg_list.push(new_bg)
	bg_list = bg_list.join("^")
	localStorage.setItem("sis50-background", bg_list)
	updateImgList()
}

function removeBackground(event) {
	var bg_to_remove = document.getElementById("bg-confirm")
			.getAttribute("data-selected")

	var bg_list = localStorage.getItem("sis50-background").split("^")

	var index = bg_list.indexOf(bg_to_remove);
	if (index > -1) {
		bg_list.splice(index, 1);
	}
	bg_list.
	bg_list = bg_list.join("^")
	localStorage.setItem("sis50-background", bg_list)
	updateImgList()

}
