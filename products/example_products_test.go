package products

import (
	"fmt"
	"gifthub/tests"
)

func ExampleSerializeMeta() {
	ctx := tests.Context()
	m := map[string]string{
		"color": "blue",
	}
	fmt.Println(SerializeMeta(ctx, m, ";"))
	// Output: color_blue
}

func ExampleUnSerializeMeta() {
	ctx := tests.Context()
	fmt.Println(UnSerializeMeta(ctx, "color_blue;size_l", ";"))
	// Output: map[color:blue size:l]
}
