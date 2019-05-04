package loglineparser

import (
	"strconv"
	"strings"
	"unicode"
)

/*
IsBlank checks if a string is whitespace or empty (""). Observe the following behavior:
    IsBlank("")        = true
    IsBlank(" ")       = true
    IsBlank("bob")     = false
    IsBlank("  bob  ") = false
Parameter:
    str - the string to check
Returns:
    true - if the string is whitespace or empty ("")
*/
func IsBlank(str string) bool {
	if str == "" {
		return true
	}

	for _, s := range str {
		if !unicode.IsSpace(s) {
			return false
		}
	}

	return true
}

func IsNumeric(s string) bool {
	for _, r := range s {
		if !(r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}

func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}

func SplitN(s, sep string, trimSpace, ignoreEmpty bool) []string {
	parts := strings.SplitN(s, sep, -1)
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		if trimSpace {
			p = strings.TrimSpace(p)
		}

		if ignoreEmpty && p == "" {
			continue
		}

		result = append(result, p)
	}

	return result
}

// Split2 将s按分隔符sep分成x份，取第x份，取第1、2、...份
func Split2(s, sep string) (s0, s1 string) {
	parts := SplitN(s, sep, true, true)
	l := len(parts)

	if l > 0 {
		s0 = parts[0]
	}
	if l > 1 {
		s1 = parts[1]
	}

	return s0, s1
}

// ParseInt parse s as int or return defaultValue
func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	return defaultValue
}
