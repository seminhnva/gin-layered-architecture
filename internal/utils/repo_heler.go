package utils

import "strings"

func NormalizeSortBy(sortBy string) string {
	switch sortBy {
	case "name":
		return "name"
	case "email":
		return "email"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}

func NormalizeOrder(order string) string {
	switch strings.ToUpper(order) {
	case "ASC":
		return "ASC"
	case "DESC":
		return "DESC"
	default:
		return "DESC"
	}
}
