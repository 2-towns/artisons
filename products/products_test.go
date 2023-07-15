package products

import "testing"

// TestImagePath build a path for a product image
func TestImagePath(t *testing.T) {
	pid := "123"
	index := 1
	_, p := ImagePath(pid, index)
	expected := "../web/images/123/1"

	if p != expected {
		t.Fatal(`The import should failed because the csv header is not correct`)
	}
}
