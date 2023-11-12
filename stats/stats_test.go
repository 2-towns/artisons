package stats

import (
	"gifthub/tests"
	"testing"
)

func TestVisitReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := Visit(c); err != nil {
		t.Fatalf(`Visit(c,"test") = %s, want nil`, err.Error())
	}
}

func TestVisitsReturnsTheVisitsForOneMonthWhenSuccess(t *testing.T) {
	c := tests.Context()

	visits, err := Visits(c)

	if err != nil {
		t.Fatalf(`Visits(c) = %v %v, want []int64{}, nil`, visits, err.Error())
	}

	if visits[0] == 0 {
		t.Fatalf(`visits[0] = 0, want > 0`)
	}

	if visits[29] == 0 {
		t.Fatalf(`visits[0] = 0, want > 0`)
	}
}

func TestUniqueVisitsReturnsTheVisitsForOneMonthWhenSuccess(t *testing.T) {
	c := tests.Context()

	visits, err := UniqueVisits(c)

	if err != nil {
		t.Fatalf(`Visits(c) = %v %v, want []int64{}, nil`, visits, err.Error())
	}

	if visits[0] == 0 {
		t.Fatalf(`visits[0] = 0, want > 0`)
	}

	if visits[29] == 0 {
		t.Fatalf(`visits[0] = 0, want > 0`)
	}
}

func TestUserReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := NewUser(c, 1); err != nil {
		t.Fatalf(`User(c, 1) = %v, want nil`, err.Error())
	}
}

func TestUsersReturnsTheUsersForOneMonthWhenSuccess(t *testing.T) {
	c := tests.Context()

	users, err := Users(c, 30)

	if err != nil {
		t.Fatalf(`Users(c, 30) = %v %v, want []int64{}, nil`, users, err.Error())
	}

	if users[0] == 0 {
		t.Fatalf(`users[0] = 0, want > 0`)
	}

	if users[29] == 0 {
		t.Fatalf(`users[0] = 0, want > 0`)
	}
}

func TestOrderReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := Order(c, "test"); err != nil {
		t.Fatalf(`Order(c, "test") = %v, want nil`, err.Error())
	}
}

func TestOrdersReturnsTheOrdersForOneMonthWhenSuccess(t *testing.T) {
	c := tests.Context()

	orders, err := Orders(c, 30)

	if err != nil {
		t.Fatalf(`Orders(c, 30) = %v %v, want []int64{}, nil`, orders, err.Error())
	}

	if orders[0] == 0 {
		t.Fatalf(`orders[0] = 0, want > 0`)
	}

	if orders[29] == 0 {
		t.Fatalf(`orders[0] = 0, want > 0`)
	}
}

func TestRevenueReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := Revenue(c, "test", 100.42); err != nil {
		t.Fatalf(`Revenue(c, 100.42) = %v, want nil`, err.Error())
	}
}

func TestRevenuesReturnsTheRevenuesForOneMonthWhenSuccess(t *testing.T) {
	c := tests.Context()

	revenues, err := Revenues(c, 30)

	if err != nil {
		t.Fatalf(`Revenues(c) = %v %v, want []int64{}, nil`, revenues, err.Error())
	}

	if revenues[0] == 0 {
		t.Fatalf(`revenues[0] = 0, want > 0`)
	}

	if revenues[29] == 0 {
		t.Fatalf(`revenues[0] = 0, want > 0`)
	}
}

func TestProductReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := SoldProduct(c, "test", "test", 1); err != nil {
		t.Fatalf(`Product(c,  "test") = %v, want nil`, err.Error())
	}
}

func TestProductsReturnsTheOrdersForOneMonthWhenSuccess(t *testing.T) {
	c := tests.Context()

	products, err := SoldProducts(c, 30)

	if err != nil {
		t.Fatalf(`Products(c , 30) = %v %v, want []int64{}, nil`, products, err.Error())
	}

	if products[0] == 0 {
		t.Fatalf(`Products[0] = 0, want > 0`)
	}

	if products[29] == 0 {
		t.Fatalf(`Products[0] = 0, want > 0`)
	}
}

func TestActiveUserReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := ActiveUser(c, "test"); err != nil {
		t.Fatalf(`ActiveUser(c, "test") = %v, want nil`, err.Error())
	}
}

func TestActiveUsersReturnsTheUsersForOneMonthWhenSuccess(t *testing.T) {
	c := tests.Context()

	users, err := ActiveUsers(c, 30)

	if err != nil {
		t.Fatalf(`ActiveUsers(c, 30) = %v %v, want []int64{}, nil`, users, err.Error())
	}

	if users[0] == 0 {
		t.Fatalf(`users[0] = 0, want > 0`)
	}

	if users[29] == 0 {
		t.Fatalf(`users[0] = 0, want > 0`)
	}
}

func TestProductViewReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := ProductView(c, "test"); err != nil {
		t.Fatalf(`ProductView(c, "test") = %v, want nil`, err.Error())
	}
}

func TestProductViewsReturnsTheUsersForOneMonthWhenSuccess(t *testing.T) {
	c := tests.Context()

	views, err := ProductViews(c, 30)

	if err != nil {
		t.Fatalf(`ProductViews(c, 30) = %v %v, want []int64{}, nil`, views, err.Error())
	}

	if len(views) == 0 {
		t.Fatalf(`len(views) = %d, want > 0, nil`, len(views))
	}

	l := len(views)
	for i := 0; i < l-1; i++ {
		if views[i].Value < views[i+1].Value {
			t.Fatalf(`views[i].Value (%d) < views[i+1].Value (%d), want views[i].Value > views[i+1].Value`, views[i].Value, views[i+1].Value)
		}
	}
}
