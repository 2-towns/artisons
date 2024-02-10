package filters

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

var filter Filter = Filter{
	Key:       "colors",
	Editable:  true,
	Label:     "colors",
	Score:     1,
	Values:    []string{"yellow", "blue"},
	Active:    true,
	UpdatedAt: time.Now(),
}
var cur string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	cur = path.Dir(filename) + "/"
}

func TestValidate(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct {
		name  string
		field string
		value string
		err   error
	}{
		{"key=", "Key", "", errors.New("input:key")},
		{"key=hello!", "Key", "hello!", errors.New("input:key")},
		{"label=", "Label", "", errors.New("input:label")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := filter

			reflect.ValueOf(&f).Elem().FieldByName(tt.field).SetString(tt.value)

			if err := f.Validate(ctx); fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestFind(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/filters.redis")

	var tests = []struct {
		name string
		key  string
		err  error
	}{
		{"key=colors", "colors", nil},
		{"key=idontexist", "idontexist", errors.New("oops the data is not found")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Find(ctx, tt.key); fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestSave(t *testing.T) {
	c := tests.Context()

	if _, err := filter.Save(c); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
}

func TestList(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/filters.redis")

	r, err := List(ctx, 0, 10)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	if r.Total == 0 {
		t.Fatalf(`total = %d, want > 0`, r.Total)
	}

	if len(r.Filters) == 0 {
		t.Fatalf(`len(filters) = %d, want > 0`, len(r.Filters))
	}

	filter := r.Filters[0]

	if filter.Key == "" {
		t.Fatalf(`key = %s, want not empty`, filter.Key)
	}
}

func TestActives(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/filters.redis")

	filters, err := Actives(ctx)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	if len(filters) == 0 {
		t.Fatalf(`len(filters) = %d, want > 0`, len(filters))
	}

	filter := filters[0]

	if filter.Key == "" {
		t.Fatalf(`key = %s, want not empty`, filter.Key)
	}
}

func TestDelete(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/filters.redis")

	var tests = []struct {
		name string
		key  string
		err  error
	}{
		{"key=", "", errors.New("input:key")},
		{"key=filters", "filters", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(ctx, tt.key); fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestExists(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/filters.redis")

	var tests = []struct {
		name   string
		key    string
		exists bool
	}{
		{"key=colors", "colors", true},
		{"key=idontexist", "idontexist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if exists, err := Exists(ctx, tt.key); err != nil || exists != tt.exists {
				t.Fatalf("exists = %v, err = %v, want %v, nil", exists, err, tt.exists)
			}
		})
	}
}

func TestEditable(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/filters.redis")

	var tests = []struct {
		name   string
		key    string
		exists bool
	}{
		{"key=colors", "colors", false},
		{"key=sizes", "sizes", true},
		{"key=idontexist", "idontexist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if exists, err := Editable(ctx, tt.key); err != nil || exists != tt.exists {
				t.Fatalf("exists = %v, err = %v, want %v, nil", exists, err, tt.exists)
			}
		})
	}
}
