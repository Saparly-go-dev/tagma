package handler

import (
	"net/http"
	"testing"
)

func TestFurniturePointCreateSupportsLegacyGet(t *testing.T) {
	routes := NewHandler(nil).InitRoutes().Routes()

	for _, route := range routes {
		if route.Method == http.MethodGet && route.Path == "/api/furniture/point_create" {
			return
		}
	}

	t.Fatal("legacy GET /api/furniture/point_create route is not registered")
}

func TestPaymentAddDebtSupportsLegacyGet(t *testing.T) {
	routes := NewHandler(nil).InitRoutes().Routes()

	for _, route := range routes {
		if route.Method == http.MethodGet && route.Path == "/api/payments/add_debt" {
			return
		}
	}

	t.Fatal("legacy GET /api/payments/add_debt route is not registered")
}

func TestCanMutateLimitsOperationalRoles(t *testing.T) {
	tests := []struct {
		role string
		path string
		want bool
	}{
		{role: "agent", path: "/api/orders/save", want: true},
		{role: "agent", path: "/api/agents/1/lock", want: false},
		{role: "ekspeditor", path: "/api/ekspeditors/product/delivered", want: true},
		{role: "ekspeditor", path: "/api/products/1", want: false},
		{role: "viewer", path: "/api/products/1", want: false},
	}

	for _, test := range tests {
		if got := canMutate(test.role, test.path); got != test.want {
			t.Errorf("canMutate(%q, %q) = %v, want %v", test.role, test.path, got, test.want)
		}
	}
}
