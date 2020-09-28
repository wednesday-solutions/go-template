package convert

import (
	"github.com/volatiletech/null"
	"strconv"
)

// StringToInt ...
func StringToInt(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}

// StringToBool ...
func StringToBool(v string) bool {
	i, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return i
}

// NullDotStringToPointerString ...
func NullDotStringToPointerString(v null.String) *string {
	return v.Ptr()
}

// PointerStringToNullDotInt ...
func PointerStringToNullDotInt(s *string) null.Int {
	if s == nil {
		return null.IntFrom(0)
	}
	v := *s
	i, err := strconv.Atoi(v)
	if err != nil {
		return null.IntFrom(0)
	}
	return null.IntFrom(i)
}
