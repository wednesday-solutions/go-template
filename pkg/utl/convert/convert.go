package convert

import (
	"strconv"

	"github.com/volatiletech/null/v8"
)

// StringToPointerString returns pointer string value
func StringToPointerString(v string) *string {
	return &v
}

// StringToInt converts string to integer
func StringToInt(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}

// StringToBool converts string to boolean
func StringToBool(v string) bool {
	i, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return i
}

// NullDotStringToPointerString converts nullable string to its pointer value
func NullDotStringToPointerString(v null.String) *string {
	return v.Ptr()
}

// NullDotStringToString converts nullable string to its value
func NullDotStringToString(v null.String) string {
	if v.Ptr() == nil {
		return ""
	}
	return *v.Ptr()
}

// NullDotIntToInt converts nullable int to its value
func NullDotIntToInt(v null.Int) int {
	if v.Ptr() == nil {
		return 0
	}
	return *v.Ptr()
}

// NullDotBoolToPointerBool converts nullable boolean to its pointer value
func NullDotBoolToPointerBool(v null.Bool) *bool {
	return v.Ptr()
}

// PointerStringToNullDotInt converts pointer string to nullable integer if present else returns default nullable value
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

func NullDotTimeToPointerInt(t null.Time) *int {
	var i int
	if t.Valid {
		i = int(t.Time.UnixMilli())
		return &i
	}
	return nil
}
