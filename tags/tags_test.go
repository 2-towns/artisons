package tags

import (
	"artisons/tests"
	"errors"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"
	"time"
)

var tag Tag = Tag{
	Key:       "phones",
	Label:     "Phones",
	Root:      false,
	UpdatedAt: time.Now(),
}

var cur string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	cur = path.Dir(filename) + "/"
}

func TestValidate(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct{ name, field, value, want string }{
		{"key=", "Key", "", "input:key"},
		{"key=!!!", "Key", "", "input:key"},
		{"label=", "Label", "", "input:label"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := tag

			reflect.ValueOf(&ta).Elem().FieldByName(tt.field).SetString(tt.value)

			if err := ta.Validate(ctx); err == nil || err.Error() != tt.want {
				t.Fatalf(`err = %v, want %s`, err, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	ctx := tests.Context()

	var cases = []struct {
		name string
		tag  string
		err  error
	}{
		{"tag=mens", "mens", nil},
		{"tag=idontexist", "idontexist", errors.New("oops the data is not found")},
	}

	tests.ImportData(ctx, cur+"testdata/tags.redis")

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			tag, err := Find(ctx, tt.tag)

			if fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want nil`, err)
			}

			if tt.err != nil {
				return
			}

			if tag.Key == "" {
				t.Fatalf(`key = %s, want not empty`, tag.Key)
			}

			if tag.Label == "" {
				t.Fatalf(`label = %s, want not empty`, tag.Label)
			}
		})
	}
}

func TestSave(t *testing.T) {
	c := tests.Context()

	if _, err := tag.Save(c); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
}

func TestList(t *testing.T) {
	c := tests.Context()

	r, err := List(c, 0, 10)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	if r.Total == 0 {
		t.Fatalf(`total = %d, want > 0`, r.Total)
	}

	if len(r.Tags) == 0 {
		t.Fatalf(`len(tags) = %d, want > 0`, len(r.Tags))
	}

	tag := r.Tags[0]

	if tag.Key == "" {
		t.Fatalf(`key = %s, want not empty`, tag.Key)
	}
}

func TestDelete(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/tags.redis")

	var cases = []struct {
		name string
		tag  string
		err  error
	}{
		{"tag=", "", errors.New("input:key")},
		{"tag=children", "children", nil},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(ctx, tt.tag); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %s`, err, tt.err)
			}
		})
	}
}

func TestExists(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/tags.redis")

	var cases = []struct {
		name string
		tag  string
		err  error
	}{
		{"tag=mens", "mens", nil},
		{"tag=idontexist", "idontexist", nil},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Exists(ctx, tt.tag); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %s`, err, tt.err)
			}
		})
	}
}

func TestElligible(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/tags.redis")

	var cases = []struct {
		name     string
		keys     []string
		eligible bool
	}{
		{"elligible", []string{"clothes", "shoes"}, true},
		{"not elligibile", []string{"mens"}, false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if eligible, err := AreEligible(ctx, tt.keys); eligible != tt.eligible || err != nil {
				t.Fatalf(`eligible = %v, err = %v, want %v, nil`, eligible, tt.eligible, err)
			}
		})
	}
}
