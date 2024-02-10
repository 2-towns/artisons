package tree

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

func TestBuild(t *testing.T) {
	ctx := tests.Context()

	tests.Del(ctx, "tag")
	tests.ImportData(ctx, cur+"testdata/tree.redis")

	tree, err := Build(ctx)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	if len(tree) != 2 {
		t.Fatalf(`len(tree) = %v, want 2`, len(tree))
	}

	mens := tree[0]

	if mens.Key != "mens" {
		t.Fatalf(`mens.Name = %v, want '%s'`, mens.Key, "mens")
	}

	if len(mens.Branches) != 2 {
		t.Fatalf(`len(mens.Branches) = %v, want 2`, len(mens.Branches))
	}

	if mens.Branches[0].Key != "clothes" {
		t.Fatalf(`mens.Branches[0].Name = %v, want '%s'`, mens.Branches[0].Key, "clothes")
	}

	if mens.Branches[1].Key != "shoes" {
		t.Fatalf(`mens.Branches[0].Name = %v, want '%s'`, mens.Branches[1].Key, "shoes")
	}

	womens := tree[1]

	if womens.Key != "womens" {
		t.Fatalf(`womens.Name = %v, want 'womens'`, womens.Key)
	}

	if len(womens.Branches) != 2 {
		t.Fatalf(`len(womens.Branches) = %v, want 2`, len(womens.Branches))
	}

	if womens.Branches[0].Key != "clothes" {
		t.Fatalf(`womens.Branches[0].Name = %v, want '%s'`, womens.Branches[0].Key, "clothes")
	}

	if womens.Branches[1].Key != "shoes" {
		t.Fatalf(`womens.Branches[0].Name = %v, want '%s'`, womens.Branches[1].Key, "shoes")
	}
}
