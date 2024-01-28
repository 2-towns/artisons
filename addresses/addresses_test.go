package addresses

import (
	"artisons/tests"
	"testing"
)

func TestGetReturnsItemsWhenAddressExists(t *testing.T) {
	ctx := tests.Context()

	res, err := Get(ctx, "8 bd du port", 2)

	if err != nil || len(res) != 2 {
		t.Fatalf(`Get(ctx, "8 bd du port", 2) = %v, %v, want not empty, nil`, res, err.Error())
	}
}
