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
	var elements = document.querySelectorAll(".tags")
	elements.forEach(element => {
		console.info(element.id)

		var config = {};

		if (element.classList.contains("tags-create")) {
			config.create = true
		}

		new TomSelect("#" + element.id, config);
	})

}

