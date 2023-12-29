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
