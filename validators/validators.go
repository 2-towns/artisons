package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var V = validator.New()

func init() {
	V.RegisterValidation("title", title)
}

// title validates a title product by allowing only necessary chars.
func title(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9()\- ]+$`).MatchString(fl.Field().String())
}
