document.body.addEventListener("otp", function () {
    var inputs = document.querySelectorAll('[data-code-input]');
    console.info(inputs)
    for (let i = 0; i < inputs.length; i++) {
        inputs[i].addEventListener('input', function (e) {
            // If the input field has a character, and there is a next input field, focus it
            if (e.target.value.length === e.target.maxLength && i + 1 < inputs.length) {
                inputs[i + 1].focus();
            }
        });

        inputs[i].addEventListener('keydown', function (e) {
            // If the input field is empty and the keyCode for Backspace (8) is detected, and there is a previous input field, focus it
            if (e.target.value.length === 0 && e.keyCode === 8 && i > 0) {
                inputs[i - 1].focus();
            }
        });

        inputs[i].addEventListener('paste', function (e) {
            const paste = (event.clipboardData || window.clipboardData).getData("text");
            const value = parseInt(paste, 10)

            if (isNaN(value)) {
                return
            }

            for (let i = 0; i < paste.length; i++) {
                if (!inputs[i]) {
                    break
                }

                inputs[i].value = paste[i]
            }

            inputs[paste.length - 1].focus();
        });
    }

    console.debug("otp event registered")
})

