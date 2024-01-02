package images

import (
	"fmt"
)

func ExampleURL() {
	o := Options{
		Width:       "60",
		Height:      "60",
		Cachebuster: 123,
	}

	fmt.Println(URL("PDT1.jpg", o))
	// Output: http://localhost:8000/resize:fill:60:60/cachebuster:123/plain/local:///PDT1.jpg

}
