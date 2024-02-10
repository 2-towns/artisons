package addresses

import (
	"artisons/tests"
	"path"
	"runtime"
	"testing"
)

var cur string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	cur = path.Dir(filename) + "/"
}

func TestGet(t *testing.T) {
	ctx := tests.Context()

	res, err := Get(ctx, "8 bd du port", 2)

	if err != nil {
		t.Fatalf(`err = %v, want not nil`, err.Error())
	}

	if len(res) != 2 {
		t.Fatalf(`len(res) =  %d, want 2`, len(res))
	}
}
