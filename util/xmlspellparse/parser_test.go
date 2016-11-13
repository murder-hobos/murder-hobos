package xmlspellparse

import (
	"database/sql"
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
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

func TestDbSpell_FromXMLSpell(t *testing.T) {
	type fields struct {
		Name          string
		Level         string
		School        string
		CastTime      string
		Duration      string
		Range         string
		Components    *Components
		Concentration bool
		Ritual        bool
		Description   string
		SourceID      int
	}
	type args struct {
		x *XMLSpell
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Zone of Truth. This should work.",
			fields{
				Name:     "Zone of Truth",
				Level:    "2",
				School:   "Enchantment",
				CastTime: "1 action",
				Duration: "10 minutes",
				Range:    "60 feet",
				Components: &Components{
					Verb: true,
					Som:  true,
					Mat:  true,
					Matdesc: sql.NullString{
						String: "",
						Valid:  false,
					},
				},
				Concentration: false,
				Ritual:        false,
				Description:   "&lt;p&gt;You create a magical zone that guards against deception in a 15-foot-radius sphere centered on a point of your choice within range. Until the spell ends, a creature that enters the spell&#39;s area for the first time on a turn or starts its turn there must make a Charisma saving throw. On a failed save, a creature can&#39;t speak a deliberate lie while in the radius. You know whether each creature succeeds or fails on its saving throw.&lt;/p&gt;&lt;p&gt;An affected creature is aware of the spell and can thus avoid answering questions to which it would normally respond with a lie. Such creatures can be evasive in its answers as long as it remains within the boundaries of the truth.&lt;/p&gt;",
				SourceID:      1,
			},
			args{
				&XMLSpell{
					XMLName:    xml.Name{"", "spell"},
					Name:       "Zone of Truth",
					Level:      "2",
					School:     "EN",
					Ritual:     "",
					Time:       "1 action",
					Range:      "60 feet",
					Components: "V, S, M",
					Duration:   "10 minutes",
					Classes:    "Bard, Cleric, Paladin, Paladin (Devotion), Paladin (Crown)",
					Texts: []string{"You create a magical zone that guards against deception in a 15-foot-radius sphere centered on a point of your choice within range. Until the spell ends, a creature that enters the spell's area for the first time on a turn or starts its turn there must make a Charisma saving throw. On a failed save, a creature can't speak a deliberate lie while in the radius. You know whether each creature succeeds or fails on its saving throw.",
						"",
						"An affected creature is aware of the spell and can thus avoid answering questions to which it would normally respond with a lie. Such creatures can be evasive in its answers as long as it remains within the boundaries of the truth.",
					},
				},
			},
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		d := &DbSpell{
			Name:          tt.fields.Name,
			Level:         tt.fields.Level,
			School:        tt.fields.School,
			CastTime:      tt.fields.CastTime,
			Duration:      tt.fields.Duration,
			Range:         tt.fields.Range,
			Components:    tt.fields.Components,
			Concentration: tt.fields.Concentration,
			Ritual:        tt.fields.Ritual,
			Description:   tt.fields.Description,
			SourceID:      tt.fields.SourceID,
		}
		var got DbSpell
		got.FromXMLSpell(tt.args.x)
		if !reflect.DeepEqual(&got, d) {
			spew.Dump(&got)
			spew.Dump(d)
			t.Errorf("%q. DbSpell.fromXMLSpell() = %v, want %v", tt.name, got, d)
		}
		if err := d.FromXMLSpell(tt.args.x); (err != nil) != tt.wantErr {
			t.Errorf("%q. DbSpell.fromXMLSpell() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_surround(t *testing.T) {
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
		if got := surround(tt.args.original, tt.args.start, tt.args.end); got != tt.want {
			t.Errorf("%q. Surround() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
