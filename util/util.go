// Package util exports a few basic utilties
package util

import (
	"bytes"
	"database/sql"
	"strconv"
	"strings"
)

// ToNullString converts a regular string to a sql.NullString
func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// ToNullInt64 converts a value to a sql.NullInt64
// different types of intergers, as well as strings representing
// integers are valid arguments
func ToNullInt64(i interface{}) sql.NullInt64 {
	switch val := i.(type) {
	case int:
	case int8:
	case int16:
	case int32:
		return sql.NullInt64{Int64: int64(val), Valid: true}
	case int64:
		return sql.NullInt64{Int64: val, Valid: true}
	case string:
		v, err := strconv.ParseInt(val, 0, 64)
		if err != nil {
			return sql.NullInt64{Int64: 0, Valid: false}
		}
		return sql.NullInt64{Int64: v, Valid: true}
	}
	return sql.NullInt64{Int64: 0, Valid: false}
}

// CapitalizeAtIndex capitalizes a single char from a string at specified index
// If an error is encountered (normally index being out of range),
// ok will be set to false and the original string returned unaltered
func CapitalizeAtIndex(s string, i int) (string, bool) {
	if i < 0 || i > len(s)-1 {
		return s, false
	}
	// TODO: Fix this ugly inefficient crap
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

// Intersperse inserts `inter` inbetween each character of `original`
// `inter` is not placed at the beginning nor end of `original`
func Intersperse(original, intr string) string {
	if original == "" || intr == "" || len(original) == 1 {
		return original
	}
	b := &bytes.Buffer{}
	for _, c := range original {
		b.WriteRune(c)
		b.WriteString(intr)
	}
	bs := b.Bytes()
	// trim off last intr
	return string(bytes.TrimRight(bs, intr))
}

// URLToCapsAndSpaces (horrible name, I just can't think of
// a more descriptive one) converts a string of the form
// `acid-splash` to `Acid Splash`
func URLToCapsAndSpaces(s string) string {
	spaces := strings.Replace(s, "-", " ", -1)
	return strings.Title(spaces)
}

// FormatURL converts a string of the form
// `Acid Splash` to `acid-splash`
func FormatURL(s string) string {
	lower := strings.ToLower(s)
	return strings.Replace(lower, " ", "-", -1)
}
