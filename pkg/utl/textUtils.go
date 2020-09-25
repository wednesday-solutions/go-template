package utl

import "github.com/volatiletech/null"

// FromNullableString ...
func FromNullableString(str null.String) string {
	if str.Valid {
		return str.String
	}
	return ""
}

// FromNullableInt ...
func FromNullableInt(str null.Int) int {
	if str.Valid {
		return str.Int
	}
	return 1
}
