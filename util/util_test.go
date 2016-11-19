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
			"Start, not end",
			args{"spell description", "first", ""},
			"firstspell description",
		},
		{
			"No start/end",
			args{"textextext", "", ""},
			"textextext",
		},
		{
			"end, not start",
			args{"basetext", "", "end"},
			"basetextend",
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

func TestToNullString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want sql.NullString
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := ToNullString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ToNullString() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCapitalizeAtIndex(t *testing.T) {
	type args struct {
		s string
		i int
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, got1 := CapitalizeAtIndex(tt.args.s, tt.args.i)
		if got != tt.want {
			t.Errorf("%q. CapitalizeAtIndex() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. CapitalizeAtIndex() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}

func TestSurround(t *testing.T) {
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := Surround(tt.args.original, tt.args.start, tt.args.end); got != tt.want {
			t.Errorf("%q. Surround() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestIntersperse(t *testing.T) {
	type args struct {
		original string
		intr     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"1 char",
			args{"V", ", "},
			"V",
		},
		{
			"2 char",
			args{"VS", ", "},
			"V, S",
		},
		{
			"3 char",
			args{"VSM", ", "},
			"V, S, M",
		},
		{
			"Empty original",
			args{"", ", "},
			"",
		},
	}
	for _, tt := range tests {
		if got := Intersperse(tt.args.original, tt.args.intr); got != tt.want {
			t.Errorf("%q. Intersperse() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
