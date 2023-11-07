package stringutil

import (
	"fmt"
)

func ExampleSlugify() {
	fmt.Println(Slugify("VERy nice title 12"))
	// Output: very-nice-title-12
}
