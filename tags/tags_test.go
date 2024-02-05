package tags

import (
	"artisons/tests"
	"testing"
	"time"
)

var tag Tag = Tag{
	Key:       "phones",
	Label:     "Phones",
	Root:      false,
	UpdatedAt: time.Now(),
}

func TestValidateReturnsErrorWhenTheKeyIsEmpty(t *testing.T) {
	c := tests.Context()

	ta := tag
	ta.Key = ""

	if err := ta.Validate(c); err == nil || err.Error() != "input:key" {
		t.Fatalf(`ta.Validate(c) = %v, want 'input:key'`, err.Error())
	}
}

func TestValidateReturnsErrorWhenTheKeyIsInvalid(t *testing.T) {
	c := tests.Context()

	ta := tag
	ta.Key = "!!!!"

	if err := ta.Validate(c); err == nil || err.Error() != "input:key" {
		t.Fatalf(`ta.Validate(c) = %v, want 'input:key'`, err.Error())
	}
}

func TestValidateReturnsErrorWhenTheLabelIsEmpty(t *testing.T) {
	c := tests.Context()

	ta := tag
	ta.Label = ""

	if err := ta.Validate(c); err == nil || err.Error() != "input:label" {
		t.Fatalf(`ta.Validate(c) = %v, want 'input:label'`, err.Error())
	}
}

func TestFindReturnsTagWhenTheKeyExists(t *testing.T) {
	c := tests.Context()

	tag, err := Find(c, tests.Tag)

	if err != nil {
		t.Fatalf(`Find(c, tests.Tag) = %v, want nil`, err.Error())
	}

	if tag.Key == "" {
		t.Fatalf(`tag.Key = %s, want not empty`, tag.Key)
	}

	if tag.Label == "" {
		t.Fatalf(`tag.Label = %s, want not empty`, tag.Label)
	}
}

func TestFindReturnsEmptyTagWhenTheKeyDoesNotExist(t *testing.T) {
	c := tests.Context()

	if _, err := Find(c, tests.DoesNotExist); err == nil || err.Error() != "oops the data is not found" {
		t.Fatalf(`Find(c, tests.DoesNotExist) = %v, want nil`, err.Error())
	}
}

func TestSaveReturnsNilWhenEmptyWhenSuccess(t *testing.T) {
	c := tests.Context()

	if _, err := tag.Save(c); err != nil {
		t.Fatalf(`tag.Save(c) = %v, want nil`, err)
	}
}

func TestListReturnsTags(t *testing.T) {
	c := tests.Context()

	r, err := List(c, 0, 10)
	if err != nil {
		t.Fatalf(`List(c, 0, 10) = %v, want nil`, err)
	}

	if r.Total == 0 {
		t.Fatalf(`r.Total = %d, want > 0`, r.Total)
	}

	if len(r.Tags) == 0 {
		t.Fatalf(`len(tags) = %d, want > 0`, len(r.Tags))
	}

	tag := r.Tags[0]

	if tag.Key == "" {
		t.Fatalf(`tag.Key = %s, want not empty`, tag.Key)
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

	if err := Delete(c, tests.TagToDelete); err != nil {
		t.Fatalf(`Delete(c, tests.TagToDelete) = %v, want nil`, err)
	}
}

func TestExistsReturnsTrueWhenTagExists(t *testing.T) {
	c := tests.Context()

	if exists, err := Exists(c, tests.Tag); exists == false || err != nil {
		t.Fatalf(`Exists(c, tests.Tag) = %v, %v, want true, nil`, exists, err)
	}
}

func TestExistsReturnsFalseWhenTagDoesNotExist(t *testing.T) {
	c := tests.Context()

	if exists, err := Exists(c, tests.DoesNotExist); exists == true || err != nil {
		t.Fatalf(`Exists(c, tests.DoesNotExist) = %v, %v, want false', nil`, exists, err)
	}
}

func TestTreeReturnsTags(t *testing.T) {
	c := tests.Context()

	tree, err := tree(c)

	if err != nil {
		t.Fatalf(`Tree(c) = %v, want nil`, err)
	}

	if len(tree) != 2 {
		t.Fatalf(`len(tree) = %v, want 2`, len(tree))
	}

	mens := tree[0]

	if mens.Key != tests.Tag {
		t.Fatalf(`mens.Name = %v, want '%s'`, mens.Key, tests.Tag)
	}

	if len(mens.Branches) != 2 {
		t.Fatalf(`len(mens.Branches) = %v, want 2`, len(mens.Branches))
	}

	if mens.Branches[0].Key != tests.Branch1 {
		t.Fatalf(`mens.Branches[0].Name = %v, want '%s'`, mens.Branches[0].Key, tests.Branch1)
	}

	if mens.Branches[1].Key != tests.Branch2 {
		t.Fatalf(`mens.Branches[0].Name = %v, want '%s'`, mens.Branches[1].Key, tests.Branch2)
	}

	womens := tree[1]

	if womens.Key != tests.Tag2 {
		t.Fatalf(`womens.Name = %v, want 'womens'`, womens.Key)
	}

	if len(womens.Branches) != 2 {
		t.Fatalf(`len(womens.Branches) = %v, want 2`, len(womens.Branches))
	}

	if womens.Branches[0].Key != tests.Branch1 {
		t.Fatalf(`womens.Branches[0].Name = %v, want '%s'`, womens.Branches[0].Key, tests.Branch1)
	}

	if womens.Branches[1].Key != tests.Branch2 {
		t.Fatalf(`womens.Branches[0].Name = %v, want '%s'`, womens.Branches[1].Key, tests.Branch2)
	}
}

func TestEligibleReturnsNilWhenTagsAreNotRooot(t *testing.T) {
	c := tests.Context()

	if eligible, err := AreEligible(c, []string{tests.Branch1, tests.Branch2}); !eligible || err != nil {
		t.Fatalf(`AreEligible(c, []string{tests.Branch1, tests.Branch2}) = %v, %v, want true, nil`, eligible, err)
	}
}

func TestEligibleReturnsErrorWhenTagsAreRoot(t *testing.T) {
	c := tests.Context()

	if eligible, err := AreEligible(c, []string{tests.Tag}); eligible || err != nil {
		t.Fatalf(`AreEligible(c, []string{tests.Tag}) = %v, %v, want false, nil`, eligible, err)
	}
}
