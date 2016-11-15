// Package util exports a few basic utilties
package util

import (
	"bytes"
	"database/sql"
	"strings"
)

// ToNullString converts a regular string to a sql.NullString
func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// CapitalizeAtIndex capitalizes a single char from a string at specified index
// If an error is encountered (normally index being out of range),
// ok will be set to false and the original string returned unaltered
func CapitalizeAtIndex(s string, i int) (string, bool) {
	if i < 0 || i > len(s)-1 {
		return s, false
	}
	out := []rune(s)
	badstr := string(out[i])
	goodstr := strings.ToUpper(badstr)
	goodrune := []rune(goodstr)
	out[i] = goodrune[0]
	return string(out), true
}

// Surround places `start` and the beginning and `end` at the end of
// an `original` string.
func Surround(original, start, end string) string {
	var b bytes.Buffer

	b.Write([]byte(start))
	b.Write([]byte(original))
	b.Write([]byte(end))

	return b.String()
}
