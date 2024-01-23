htmx.defineExtension("alert", {
	onEvent: function (name) {
		if (name === "htmx:beforeSwap") {
			document.getElementById("alert").innerHTML = "";
		}
	},
});

htmx.defineExtension("input", {
	onEvent: function (name) {
		if (name === "htmx:beforeSwap") {
			document.querySelectorAll(".form-error").forEach((e) => {
				e.remove();
			});
		}
	},
});


if (document.getElementsByClassName("tags").length > 0) {
	var config = {};
	var element = document.querySelector(".tags")
	if (element.classList.contains("tags-create")) {
		config.create = true
	}
	console.info(element, config)
	new TomSelect('.tags', config);
}

