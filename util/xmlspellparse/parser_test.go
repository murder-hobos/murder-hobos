package xmlspellparse

import (
	"database/sql"
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
)

func TestComponents_parseComponents(t *testing.T) {
	type fields struct {
		Verb    bool
		Som     bool
		Mat     bool
		Matdesc sql.NullString
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"V",
			fields{Verb: true, Som: false, Mat: false, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
			args{"V"},
		},
		{
			"V, S",
			fields{Verb: true, Som: true, Mat: false, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
			args{"V, S"},
		},
		{
			"V, S, M",
			fields{Verb: true, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
			args{"V, S, M"},
		},
		{
			"V, S, M (text)",
			fields{Verb: true, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "A jade circlet worth at least 1,500 gp, which you must place on your head before you cast the spell",
				Valid:  true},
			},
			args{"V, S, M (a jade circlet worth at least 1,500 gp, which you must place on your head before you cast the spell)"},
		},

		{
			"S",
			fields{Verb: false, Som: true, Mat: false, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
			args{"S"},
		},
		{
			"V, M (text)",
			fields{Verb: true, Som: false, Mat: true, Matdesc: sql.NullString{
				String: "Bat fur and a drop of pitch or piece of coal",
				Valid:  true},
			},
			args{"V, M (bat fur and a drop of pitch or piece of coal)"},
		},
		{
			"S, M",
			fields{Verb: false, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "",
				Valid:  false},
			},
			args{"S, M"},
		},
		{
			"S, M (text)",
			fields{Verb: false, Som: true, Mat: true, Matdesc: sql.NullString{
				String: "A glowing stick of incense or a crystal vial filled with phosphorescent material",
				Valid:  true},
			},
			args{"S, M (a glowing stick of incense or a crystal vial filled with phosphorescent material)"},
		},
	}
	for _, tt := range tests {
		c := &Components{
			Verb:    tt.fields.Verb,
			Som:     tt.fields.Som,
			Mat:     tt.fields.Mat,
			Matdesc: tt.fields.Matdesc,
		}
		c.parseComponents(tt.args.s)
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

func TestXMLSpell_ToDbSpell(t *testing.T) {
	type fields struct {
		XMLName    xml.Name
		Name       string
		Level      string
		School     string
		Ritual     string
		Time       string
		Range      string
		Components string
		Duration   string
		Classes    string
		Texts      []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    DbSpell
		wantErr bool
	}{
		{
			"Zone of truth. Should work.",
			fields{
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
			DbSpell{
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
				SourceID:      PHBid,
			},
			false,
		},
		{
			"Absorb Elements. Should be EE.",
			fields{
				Name:       "Absorb Elements (EE)",
				Level:      "1",
				School:     "A",
				Ritual:     "",
				Time:       "1 reaction, which you take when you take acid, cold, fire, lightning, or thunder damage",
				Range:      "",
				Components: "S",
				Duration:   "1 round",
				Classes:    "Druid, Ranger, Wizard, Fighter (Eldritch Knight)",
				Texts: []string{"The spell captures some of the incoming energy, lessening its effect on you and storing it for your next melee attack. You have resistance to the triggering damage type until the start of your next turn. Also, the first time you hit with a melee attack on your next turn, the target takes an extra 1d6 damage of the triggering type, and the spell ends.",
					"",
					"At Higher Levels: When you cast this spell using a spell slot of 2nd level or higher, the extra damage increases by 1d6 for each slot level above 1st.",
					"",
					"This spell can be found in the Elemental Evil Player's Companion",
				},
			},
			DbSpell{
				Name:     "Absorb Elements (EE)",
				Level:    "1",
				School:   "Abjuration",
				CastTime: "1 reaction, which you take when you take acid, cold, fire, lightning, or thunder damage",
				Duration: "1 round",
				Range:    "",
				Components: &Components{
					Verb: false,
					Som:  true,
					Mat:  false,
					Matdesc: sql.NullString{
						String: "",
						Valid:  false,
					},
				},
				Concentration: false,
				Ritual:        false,
				Description:   "&lt;p&gt;The spell captures some of the incoming energy, lessening its effect on you and storing it for your next melee attack. You have resistance to the triggering damage type until the start of your next turn. Also, the first time you hit with a melee attack on your next turn, the target takes an extra 1d6 damage of the triggering type, and the spell ends.&lt;/p&gt;&lt;p&gt;At Higher Levels: When you cast this spell using a spell slot of 2nd level or higher, the extra damage increases by 1d6 for each slot level above 1st.&lt;/p&gt;&lt;p&gt;This spell can be found in the Elemental Evil Player&#39;s Companion&lt;/p&gt;",
				SourceID:      EEid,
			},
			false,
		},
	}
	for _, tt := range tests {
		x := XMLSpell{
			Name:       tt.fields.Name,
			Level:      tt.fields.Level,
			School:     tt.fields.School,
			Ritual:     tt.fields.Ritual,
			Time:       tt.fields.Time,
			Range:      tt.fields.Range,
			Components: tt.fields.Components,
			Duration:   tt.fields.Duration,
			Classes:    tt.fields.Classes,
			Texts:      tt.fields.Texts,
		}
		got, err := x.ToDbSpell()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. XMLSpell.ToDbSpell() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			spew.Dump(got, tt.want)
			t.Errorf("%q. XMLSpell.ToDbSpell() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestXMLSpell_ParseClasses(t *testing.T) {
	type fields struct {
		Name       string
		Level      string
		School     string
		Ritual     string
		Time       string
		Range      string
		Components string
		Duration   string
		Classes    string
		Texts      []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []Class
		want1  bool
	}{
		{
			"3 classes",
			fields{
				Classes: "Cleric, Cleric (Arcana), Druid",
			},
			[]Class{
				Class{
					ID:   2,
					Name: "Cleric",
					BaseClass: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				},
				Class{
					ID:   3,
					Name: "Cleric (Arcana)",
					BaseClass: sql.NullInt64{
						Int64: 2,
						Valid: true,
					},
				},
				Class{
					ID:   12,
					Name: "Druid",
					BaseClass: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				},
			},
			true,
		},
		{
			"No classes",
			fields{},
			[]Class{},
			false,
		},
		{
			"One class",
			fields{
				Classes: "Bard",
			},
			[]Class{
				Class{
					ID:   1,
					Name: "Bard",
					BaseClass: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		x := &XMLSpell{
			Name:       tt.fields.Name,
			Level:      tt.fields.Level,
			School:     tt.fields.School,
			Ritual:     tt.fields.Ritual,
			Time:       tt.fields.Time,
			Range:      tt.fields.Range,
			Components: tt.fields.Components,
			Duration:   tt.fields.Duration,
			Classes:    tt.fields.Classes,
			Texts:      tt.fields.Texts,
		}
		got, got1 := x.ParseClasses()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. XMLSpell.ParseClasses() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. XMLSpell.ParseClasses() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
