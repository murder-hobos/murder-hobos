package xmlspellparse

import (
	"database/sql"
	"reflect"
	"testing"
)

func Test_parseComponents(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *Components
	}{
		{"V",
			args{"V"},
			&Components{Verb: true, Som: false, Mat: false, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
		},
		{"V, S",
			args{"V, S"},
			&Components{Verb: true, Som: true, Mat: false, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
		},
		{"V, S, M",
			args{"V, S, M"},
			&Components{Verb: true, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
		},
		{"V, S, M (text)",
			args{"V, S, M (a jade circlet worth at least 1,500 gp, which you must place on your head before you cast the spell)"},
			&Components{Verb: true, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "A jade circlet worth at least 1,500 gp, which you must place on your head before you cast the spell",
				Valid:  true},
			},
		},
		{"S",
			args{"S"},
			&Components{Verb: false, Som: true, Mat: false, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
		},
		{"V, M (text)",
			args{"V, M (bat fur and a drop of pitch or piece of coal)"},
			&Components{Verb: true, Som: false, Mat: true, Matdesc: sql.NullString{
				String: "Bat fur and a drop of pitch or piece of coal",
				Valid:  true},
			},
		},
		{"S, M",
			args{"S, M"},
			&Components{Verb: false, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
		},
		{"S, M (text)",
			args{"S, M (a glowing stick of incense or a crystal vial filled with phosphorescent material)"},
			&Components{Verb: false, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "A glowing stick of incense or a crystal vial filled with phosphorescent material",
				Valid:  true},
			},
		},
	}
	for _, tt := range tests {
		if got := parseComponents(tt.args.s); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. parseComponents() %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_capitalizeAtIndex(t *testing.T) {
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
		if sgot, bgot := capitalizeAtIndex(tt.args.s, tt.args.i); sgot != tt.swant || bgot != tt.bwant {
			t.Errorf("%q. capitalizeAtIndex() = %v, %v, want %v, %v", tt.name, sgot, bgot, tt.swant, tt.bwant)
		}
	}
}

func Test_toNullString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want sql.NullString
	}{
		{"Should be valid", args{"valid"}, sql.NullString{"valid", true}},
		{"Should be invalid", args{""}, sql.NullString{"", false}},
	}
	for _, tt := range tests {
		if got := toNullString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. toNullString() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
