package repository

import (
	"errors"
	"strings"
)

// allowedOrderSort converts an API sort key to a trusted SQL expression.
func allowedOrderSort(value string) (string, error) {
	allowed := map[string]string{
		"id":              "orders.id",
		"id asc":          "orders.id asc",
		"id desc":         "orders.id desc",
		"created_at":      "orders.created_at",
		"created_at asc":  "orders.created_at asc",
		"created_at desc": "orders.created_at desc",
		"sum asc":         "orders.sum asc",
		"sum desc":        "orders.sum desc",
	}
	sortOrder, ok := allowed[strings.ToLower(strings.TrimSpace(value))]
	if !ok {
		return "", errors.New("invalid order sort")
	}
	return sortOrder, nil
}
