package filters

import (
	"artisons/tests"
	"testing"
	"time"
)

var filter Filter = Filter{
	Key:       tests.FilterColorKey,
	Editable:  true,
	Label:     tests.FilterColorLabel,
	Score:     1,
	Values:    []string{"yellow", "blue"},
	Active:    true,
	UpdatedAt: time.Now(),
}

func TestValidateReturnsErrorWhenTheKeyIsEmpty(t *testing.T) {
	c := tests.Context()

	f := filter
	f.Key = ""

	if err := f.Validate(c); err == nil || err.Error() != "input:key" {
		t.Fatalf(`f.Validate(c) = %v, want 'input:key'`, err.Error())
	}
}

func TestValidateReturnsErrorWhenTheKeyIsInvalid(t *testing.T) {
	c := tests.Context()

	f := filter
	f.Key = "hello!"

	if err := f.Validate(c); err == nil || err.Error() != "input:key" {
		t.Fatalf(`f.Validate(c) = %v, want 'input:key'`, err.Error())
	}
}

func TestValidateReturnsErrorWhenLabelIsEmpty(t *testing.T) {
	c := tests.Context()

	f := filter
	f.Label = ""

	if err := f.Validate(c); err == nil || err.Error() != "input:label" {
		t.Fatalf(`ta.Validate(c) = %v, want 'input:label'`, err.Error())
	}
}

func TestFindReturnsFilterWhenTheKeyExists(t *testing.T) {
	c := tests.Context()

	filter, err := Find(c, tests.FilterColorKey)

	if err != nil {
		t.Fatalf(`Find(c, tests.FilterColorKey) = %v, want nil`, err.Error())
	}

	if filter.Key == "" {
		t.Fatalf(`tag.Key = %s, want not empty`, filter.Key)
	}

	if filter.Label == "" {
		t.Fatalf(`tag.Label = %s, want not empty`, filter.Label)
	}
}

func TestFindReturnsEmptyFilterWhenTheKeyDoesNotExist(t *testing.T) {
	c := tests.Context()

	if _, err := Find(c, tests.DoesNotExist); err == nil || err.Error() != "oops the data is not found" {
		t.Fatalf(`Find(c, tests.DoesNotExist) = %v, want nil`, err.Error())
	}
}

func TestSaveReturnsNilWhenEmptyWhenSuccess(t *testing.T) {
	c := tests.Context()

	if _, err := filter.Save(c); err != nil {
		t.Fatalf(`filter.Save(c) = %v, want nil`, err)
	}
}

func TestListReturnsFilters(t *testing.T) {
	c := tests.Context()

	r, err := List(c, 0, 10)
	if err != nil {
		t.Fatalf(`List(c) = %v, want nil`, err)
	}

	if r.Total == 0 {
		t.Fatalf(`r.Total = %d, want > 0`, r.Total)
	}

	if len(r.Filters) == 0 {
		t.Fatalf(`len(r.Filters) = %d, want > 0`, len(r.Filters))
	}

	filter := r.Filters[0]

	if filter.Key == "" {
		t.Fatalf(`filter.Key = %s, want not empty`, filter.Key)
	}
}

func TestActivesReturnsFilters(t *testing.T) {
	c := tests.Context()

	filters, err := Actives(c)
	if err != nil {
		t.Fatalf(`Actives(c) = %v, want nil`, err)
	}

	if len(filters) == 0 {
		t.Fatalf(`len(filters) = %d, want > 0`, len(filters))
	}

	filter := filters[0]

	if filter.Key == "" {
		t.Fatalf(`filter.Key = %s, want not empty`, filter.Key)
	}
}

func TestDeleteReturnsErrorWhenKeyIsEmpty(t *testing.T) {
	c := tests.Context()

	if err := Delete(c, ""); err == nil || err.Error() != "input:key" {
		t.Fatalf(`Delete(c, "") = %v, want 'input:key'`, err)
	}
}

func TestDeleteReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := Delete(c, tests.FilterToDeleteKey); err != nil {
		t.Fatalf(`Delete(c, tests.FilterToDeleteKey) = %v, want nil`, err)
	}
}

func TestExistsReturnsTrueWhenFilterDoesExist(t *testing.T) {
	c := tests.Context()

	if exists, err := Exists(c, tests.FilterColorKey); !exists || err != nil {
		t.Fatalf(`Exists(c, tests.FilterColorKey) = %v, %v, want true', nil`, exists, err)
	}
}

func TestExistsReturnsFalseWhenFilterDoesNotExist(t *testing.T) {
	c := tests.Context()

	if exists, err := Exists(c, tests.DoesNotExist); exists || err != nil {
		t.Fatalf(`Exists(c, tests.DoesNotExist) = %v, %v, want false', nil`, exists, err)
	}
}

func TestEditableReturnsFalseWhenFilterIsNotEditable(t *testing.T) {
	c := tests.Context()

	if editable, err := Editable(c, tests.FilterColorKey); editable || err != nil {
		t.Fatalf(`Exists(c, tests.FilterColorKey) = %v, %v, want false', nil`, editable, err)
	}
}

func TestEditableReturnsTrueWhenFilterIsEditable(t *testing.T) {
	c := tests.Context()

	if editable, err := Editable(c, tests.FilterSizeKey); !editable || err != nil {
		t.Fatalf(`Exists(c, tests.FilterSizeKey) = %v, %v, want true', nil`, editable, err)
	}
}

func TestEditableReturnsFalseWhenEditableFilterThatDoesNotExist(t *testing.T) {
	c := tests.Context()

	if editable, err := Editable(c, tests.DoesNotExist); editable || err != nil {
		t.Fatalf(`Exists(c, tests.DoesNotExist) = %v, %v, want false', nil`, editable, err)
	}
}
