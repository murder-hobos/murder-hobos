package util

import (
	"database/sql"
	"reflect"
	"testing"
)

func Test_CapitalizeAtIndex(t *testing.T) {
	type args struct {
		s string
		i int
	}
	tests := []struct {
		name  string
		args  args
		swant string
		bwant bool
	}{
		{"First", args{"bill", 0}, "Bill", true},
		{"Negative index", args{"asdf", -1}, "asdf", false},
		{"Index too big", args{"qwert", 5}, "qwert", false},
	}
	for _, tt := range tests {
		if sgot, bgot := CapitalizeAtIndex(tt.args.s, tt.args.i); sgot != tt.swant || bgot != tt.bwant {
			t.Errorf("%q. CapitalizeAtIndex() = %v, %v, want %v, %v", tt.name, sgot, bgot, tt.swant, tt.bwant)
		}
	}
}

func Test_Surround(t *testing.T) {
	type args struct {
		original string
		start    string
		end      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Paragraph tags",
			args{"spell description", "<p>", "</p>"},
			"&lt;p&gt;spell description&lt;/p&gt;",
		},
		{
			"No start/end",
			args{"textextext", "", ""},
			"textextext",
		},
	}
	for _, tt := range tests {
		if got := Surround(tt.args.original, tt.args.start, tt.args.end); got != tt.want {
			t.Errorf("%q. Surround() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_ToNullString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want sql.NullString
	}{
		{"Should be valid", args{"valid"}, sql.NullString{String: "valid", Valid: true}},
		{"Should be invalid", args{""}, sql.NullString{String: "", Valid: false}},
	}
	for _, tt := range tests {
		if got := ToNullString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ToNullString() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
