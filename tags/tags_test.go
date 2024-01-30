package tags

import (
	"artisons/db"
	"artisons/tests"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

var tag Tag = Tag{
	Key:       "phones",
	Label:     "Phones",
	Root:      false,
	UpdatedAt: time.Now(),
}

func init() {
	ctx := tests.Context()

	db.Redis.Del(ctx, "tags")
	db.Redis.Del(ctx, "tags:root")
	db.Redis.HSet(ctx, "tag:mens", "key", "mens", "image", "tags/1.jpeg", "label", "Mens", "order", "1")
	db.Redis.HSet(ctx, "tag:womens", "key", "womens", "image", "tags/2.jpeg", "label", "Womens", "order", "2")
	db.Redis.HSet(ctx, "tag:children", "key", "children", "image", "tags/3.jpeg", "label", "Children", "order", "2")

	db.Redis.ZAdd(ctx, "tags", redis.Z{
		Score:  float64(1),
		Member: "mens",
	}, redis.Z{
		Score:  float64(2),
		Member: "womens",
	}, redis.Z{
		Score:  float64(2),
		Member: "shoes",
	}, redis.Z{
		Score:  float64(2),
		Member: "clothes",
	})

	db.Redis.ZAdd(ctx, "tags:root", redis.Z{
		Score:  float64(1),
		Member: "mens",
	}, redis.Z{
		Score:  float64(2),
		Member: "womens",
	})

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

	tag, err := Find(c, "mens")

	if err != nil {
		t.Fatalf(`Find(c, "mens") = %v, want nil`, err.Error())
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

	if _, err := Find(c, "hello"); err == nil || err.Error() != "oops the data is not found" {
		t.Fatalf(`Find(c, "hello") = %v, want nil`, err.Error())
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
		t.Fatalf(`List(c) = %v, want nil`, err)
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
		t.Fatalf(`Delete(c) = %v, want 'input:key'`, err)
	}
}

func TestDeleteReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := Delete(c, "children"); err != nil {
		t.Fatalf(`Delete(c) = %v, want nil`, err)
	}
}

func TestExistsReturnsTrueWhenTagExists(t *testing.T) {
	c := tests.Context()

	if exists, err := Exists(c, "mens"); exists == false || err != nil {
		t.Fatalf(`Exists(c, "mens") = %v, %v, want true, nil`, exists, err)
	}
}

func TestExistsReturnsFalseWhenTagDoesNotExist(t *testing.T) {
	c := tests.Context()

	if exists, err := Exists(c, "hello"); exists == true || err != nil {
		t.Fatalf(`Exists(c, "hello") = %v, %v, want false', nil`, exists, err)
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

	if mens.Key != "mens" {
		t.Fatalf(`mens.Name = %v, want 'mens'`, mens.Key)
	}

	if len(mens.Branches) != 2 {
		t.Fatalf(`len(mens.Branches) = %v, want 2`, len(mens.Branches))
	}

	if mens.Branches[0].Key != "clothes" {
		t.Fatalf(`mens.Branches[0].Name = %v, want 'clothes'`, mens.Branches[0].Key)
	}

	if mens.Branches[1].Key != "shoes" {
		t.Fatalf(`mens.Branches[0].Name = %v, want 'shoes'`, mens.Branches[1].Key)
	}

	womens := tree[1]

	if womens.Key != "womens" {
		t.Fatalf(`womens.Name = %v, want 'womens'`, womens.Key)
	}

	if len(womens.Branches) != 2 {
		t.Fatalf(`len(womens.Branches) = %v, want 2`, len(womens.Branches))
	}

	if womens.Branches[0].Key != "clothes" {
		t.Fatalf(`womens.Branches[0].Name = %v, want 'clothes'`, womens.Branches[0].Key)
	}

	if womens.Branches[1].Key != "shoes" {
		t.Fatalf(`womens.Branches[0].Name = %v, want 'shoes'`, womens.Branches[1].Key)
	}
}

func TestEligibleReturnsNilWhenTagsAreNotRooot(t *testing.T) {
	c := tests.Context()

	if eligible, err := AreEligible(c, []string{"clothes", "shoes"}); !eligible || err != nil {
		t.Fatalf(`AreEligible(c, []string{"clothes", "shoes"}) = %v, %v, want true, nil`, eligible, err)
	}
}

func TestEligibleReturnsErrorWhenTagsAreRoot(t *testing.T) {
	c := tests.Context()

	if eligible, err := AreEligible(c, []string{"mens"}); eligible || err != nil {
		t.Fatalf(`AreEligible(c, []string{"mens"}) = %v, %v, want false, nil`, eligible, err)
	}
}
