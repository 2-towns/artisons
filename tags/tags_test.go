package tags

import (
	"fmt"
	"gifthub/tests"
	"testing"
)

var tag Tag = Tag{
	Name:  "books",
	Label: "Livres",
	Score: 0,
}

func TestSaveReturnsErrorWhenEmptyWhenTheTagIsEmpty(t *testing.T) {
	c := tests.Context()

	ta := tag
	ta.Name = ""

	if err := ta.Save(c); err == nil || err.Error() != "input_tag_name_invalid" {
		t.Fatalf(`ta.Save(c) = %v, want 'input_tag_name_invalid'`, err.Error())
	}
}

func TestSaveReturnsErrorWhenEmptyWhenTheTagIsNotAlpha(t *testing.T) {
	c := tests.Context()

	ta := tag
	ta.Name = "test1"

	if err := ta.Save(c); err == nil || err.Error() != "input_tag_name_invalid" {
		t.Fatalf(`ta.Save(c) = %v, want 'input_tag_name_invalid'`, err.Error())
	}
}

func TestSaveReturnsErrorWhenEmptyLabel(t *testing.T) {
	c := tests.Context()

	ta := tag
	ta.Label = ""

	if err := ta.Save(c); err == nil || err.Error() != "input_tag_label_invalid" {
		t.Fatalf(`ta.Save(c) = %v, want 'input_tag_label_invalid'`, err.Error())
	}
}

func TestSaveReturnsNilWhenEmptyWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := tag.Save(c); err != nil {
		t.Fatalf(`tag.Save(c) = %v, want 'input_tag_name_invalid'`, err)
	}
}

func TestLinkReturnsNilWhenTagIsEmpty(t *testing.T) {
	c := tests.Context()

	if err := tag.Link(c, "", 1); err == nil || err.Error() != "tag_name_required" {
		t.Fatalf(`tag.Link(c, "", 1) = %v, want 'tag_name_required'`, err.Error())
	}
}

func TestLinkReturnsErrorWhenScoreIsZero(t *testing.T) {
	c := tests.Context()

	if err := tag.Link(c, "", 1); err == nil || err.Error() != "tag_name_required" {
		t.Fatalf(`tag.Link(c, "", 1) = %v, want 'tag_name_required'`, err.Error())
	}
}

func TestLinkReturnsErrorWhenTagIsNotFound(t *testing.T) {
	c := tests.Context()

	if err := tag.Link(c, "hello", 1); err == nil || err.Error() != "tag_not_found" {
		t.Fatalf(`tag.Link(c, "hello", 1) = %v, want 'tag_not_found'`, err.Error())
	}
}

func TestLinkReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := tag.Link(c, "arabic", 1); err != nil {
		t.Fatalf(`tag.Link(c, "arabic", 1) = %v, want nil`, err)
	}
}

func TestRemoveLinkReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	v := Tag{Name: "games"}

	if err := v.RemoveLink(c, "arabic"); err != nil {
		t.Fatalf(`Tag{Name: "games"}.RemoveLink(c, "arabic") = %v, want nil`, err)
	}
}

func TestRootReturnsRootTags(t *testing.T) {
	c := tests.Context()

	tags, err := Root(c, 3)
	if err != nil {
		t.Fatalf(`Root(c, 3) = %v, want nil`, err)
	}

	expected := "[{mens Pour les hommes 0 [{tshirts Tshirts 0 []} {books Livres 0 [{arabic Arabe 0 []}]} {clothes Vêtements 0 []}]} {womens Pour les femmes 0 [{tshirts Tshirts 0 []} {clothes Vêtements 0 []}]}]"
	if fmt.Sprintf("%v", tags) != expected {
		t.Fatalf(`tags = '%v', want [{mens Pour les hommes 0 [{tshirts Tshirts 0 []} {books Livres 0 [{arabic Arabe 0 []}]} {clothes Vêtements 0 []}]} {womens Pour les femmes 0 [{tshirts Tshirts 0 []} {clothes Vêtements 0 []}]}]`, tags)
	}
}

func TestRootWithOneDpethReturnsRootTagsForOneDepth(t *testing.T) {
	c := tests.Context()

	tags, err := Root(c, 1)
	if err != nil {
		t.Fatalf(`Root(c, 3) = %v, want nil`, err)
	}

	expected := "[{mens Pour les hommes 0 []} {womens Pour les femmes 0 []}]"
	if fmt.Sprintf("%v", tags) != expected {
		t.Fatalf(`tags = '%v', want [{mens Pour les hommes 0 []} {womens Pour les femmes 0 []}]`, tags)
	}
}

func TestListReturnsTags(t *testing.T) {
	c := tests.Context()

	tags, err := List(c)
	if err != nil {
		t.Fatalf(`List(c) = %v, want nil`, err)
	}

	if len(tags) == 0 {
		t.Fatalf(`len(tags) = %d, want > 0`, len(tags))
	}

	tag := tags[0]

	if tag.Name == "" {
		t.Fatalf(`tag.Name = %s, want not empty`, tag.Name)
	}

	if tag.Label == "" {
		t.Fatalf(`tag.Label = %s, want not empty`, tag.Label)
	}
}

func TestWithLinksReturnsTagWithItLinks(t *testing.T) {
	c := tests.Context()

	val, err := tag.WithLinks(c)
	if err != nil {
		t.Fatalf(`List(c) = %v, want nil`, err)
	}

	if len(val.Links) == 0 {
		t.Fatalf(`val.Links = %d, want > 0`, len(val.Links))
	}

	expected := "[{arabic Arabe 0 []}]"
	if fmt.Sprintf("%v", val.Links) != expected {
		t.Fatalf(`val.Links = %v, want [{arabic Arabe 0 []}]`, val.Links)
	}

	if tag.Label == "" {
		t.Fatalf(`tag.Label = %s, want not empty`, tag.Label)
	}
}
