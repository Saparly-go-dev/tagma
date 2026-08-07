package repository

import "testing"

func TestAllowedOrderSortRejectsSQLExpression(t *testing.T) {
	if _, err := allowedOrderSort("id desc; drop table users"); err == nil {
		t.Fatal("expected unsafe sort expression to be rejected")
	}
}

func TestAllowedOrderSortMapsKnownValue(t *testing.T) {
	got, err := allowedOrderSort("created_at desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "orders.created_at desc" {
		t.Fatalf("unexpected sort expression: %s", got)
	}
}
